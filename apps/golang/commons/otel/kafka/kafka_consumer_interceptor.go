package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	semconv "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type KafkaConsumer struct {
	tracer     trace.Tracer
	meter      metric.Meter
	propagator propagation.TextMapPropagator

	latency metric.Float64Histogram
}

func NewKafkaConsumer() *KafkaConsumer {

	// Instantiate trace provider
	tracer := otel.GetTracerProvider().Tracer(semconv.KafkaConsumerName)

	// Instantiate meter provider
	meter := otel.GetMeterProvider().Meter(semconv.KafkaConsumerName)

	// Instantiate propagator
	propagator := otel.GetTextMapPropagator()

	// Create HTTP client latency histogram
	latency, err := meter.Float64Histogram(
		semconv.MessagingConsumerLatencyName,
		metric.WithUnit("ms"),
		metric.WithDescription("Measures the duration of receive operation"),
		metric.WithExplicitBucketBoundaries(semconv.MessagingExplicitBucketBoundaries...),
	)
	if err != nil {
		panic(err)
	}

	return &KafkaConsumer{
		tracer:     tracer,
		meter:      meter,
		propagator: propagator,

		latency: latency,
	}
}

func (k *KafkaConsumer) Consume(
	ctx context.Context,
	msg *sarama.ConsumerMessage,
	consumerGroup string,
	consumeFunc func(context.Context) error,
) {

	// Start timer
	consumeStartTime := time.Now()

	// Get tracing info from message
	headers := propagation.MapCarrier{}

	for _, recordHeader := range msg.Headers {
		headers[string(recordHeader.Key)] = string(recordHeader.Value)
	}

	propagator := otel.GetTextMapPropagator()
	ctx = propagator.Extract(ctx, headers)

	spanAttrs := semconv.WithMessagingKafkaConsumerAttributes(msg, consumerGroup)
	attrs := semconv.WithMessagingKafkaConsumerAttributes(msg, consumerGroup)

	ctx, span := k.tracer.Start(
		ctx,
		fmt.Sprintf("%s receive", msg.Topic),
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(spanAttrs...),
	)
	defer span.End()

	// // Run the actual consume function
	// err := consumeFunc(ctx)
	// if err != nil {
	// 	attrs = append(attrs, semconv.ErrorType.String(err.Error()))
	// }

	// Record consume latency
	elapsedTime := float64(time.Since(consumeStartTime)) / float64(time.Millisecond)
	k.latency.Record(ctx, elapsedTime,
		metric.WithAttributes(
			attrs...,
		))
}
