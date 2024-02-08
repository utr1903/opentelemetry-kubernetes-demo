package kafka

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	commonerr "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/error"
	semconv "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type KafkaProducer struct {
	producer sarama.AsyncProducer

	tracer     trace.Tracer
	meter      metric.Meter
	propagator propagation.TextMapPropagator

	latency metric.Float64Histogram
}

func NewKafkaProducer(
	producer sarama.AsyncProducer,
) *KafkaProducer {

	// Instantiate trace provider
	tracer := otel.GetTracerProvider().Tracer(semconv.KafkaProducerName)

	// Instantiate meter provider
	meter := otel.GetMeterProvider().Meter(semconv.KafkaProducerName)

	// Instantiate propagator
	propagator := otel.GetTextMapPropagator()

	// Create HTTP client latency histogram
	latency, err := meter.Float64Histogram(
		semconv.MessagingProducerLatencyName,
		metric.WithUnit("ms"),
		metric.WithDescription("Measures the duration of publish operation"),
		metric.WithExplicitBucketBoundaries(semconv.MessagingExplicitBucketBoundaries...),
	)
	if err != nil {
		panic(err)
	}

	return &KafkaProducer{
		producer: producer,

		tracer:     tracer,
		meter:      meter,
		propagator: propagator,

		latency: latency,
	}
}

func (k *KafkaProducer) Publish(
	ctx context.Context,
	msg *sarama.ProducerMessage,
	errType *string,
) {

	// Start timer
	produceStartTime := time.Now()

	// Get metric attributes
	metricAttrs := semconv.WithMessagingKafkaProducerAttributes(msg)

	// Create produce latency recording function
	produceRecordFunc := func(
		ctx context.Context,
		attrs []attribute.KeyValue,
	) {
		// Record producer latency
		elapsedTime := float64(time.Since(produceStartTime)) / float64(time.Millisecond)
		k.latency.Record(ctx, elapsedTime,
			metric.WithAttributes(
				attrs...,
			))
	}

	// Inject tracing info into message
	span := k.createProducerSpan(ctx, msg)
	defer span.End()

	if errType != nil && *errType == commonerr.KAFKA_CONNECTION_ERROR {

		// Sleep as if trying to reach Kafka
		time.Sleep(3 * time.Second)

		// Create error
		err := errors.New("reaching out to kafka cluster timed out")

		// Create span attributes
		spanAttrs := []attribute.KeyValue{
			semconv.OtelStatusCode.String("ERROR"),
			semconv.OtelStatusDescription.String("Reaching out to Kafka cluster timed out."),
		}
		span.SetAttributes(spanAttrs...)
		span.RecordError(
			err,
			trace.WithAttributes(
				semconv.ExceptionEscaped.Bool(true),
			))

		// Add error as attribute and record latency
		metricAttrs = append(metricAttrs, semconv.ErrorType.String("Kafka time out"))
		produceRecordFunc(ctx, metricAttrs)
		return
	}

	// Publish message
	k.producer.Input() <- msg
	<-k.producer.Successes()

	produceRecordFunc(ctx, metricAttrs)
}

func (k *KafkaProducer) createProducerSpan(
	ctx context.Context,
	msg *sarama.ProducerMessage,
) trace.Span {
	spanAttrs := semconv.WithMessagingKafkaProducerAttributes(msg)
	spanContext, span := k.tracer.Start(
		ctx,
		fmt.Sprintf("%s publish", msg.Topic),
		trace.WithSpanKind(trace.SpanKindProducer),
		trace.WithAttributes(spanAttrs...),
	)

	carrier := propagation.MapCarrier{}
	propagator := otel.GetTextMapPropagator()
	propagator.Inject(spanContext, carrier)

	for key, value := range carrier {
		msg.Headers = append(msg.Headers, sarama.RecordHeader{Key: []byte(key), Value: []byte(value)})
	}

	return span
}
