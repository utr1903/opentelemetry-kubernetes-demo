package http

import (
	"context"
	"net/http"
	"time"

	semconv "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type Opts struct {
	Timeout time.Duration
}

type OptFunc func(*Opts)

func defaultOpts() *Opts {
	return &Opts{
		Timeout: time.Duration(30 * time.Second),
	}
}

type HttpClient struct {
	Opts *Opts

	client *http.Client

	tracer     trace.Tracer
	meter      metric.Meter
	propagator propagation.TextMapPropagator

	latency metric.Float64Histogram
}

func New(
	optFuncs ...OptFunc,
) *HttpClient {

	// Instantiate options with default values
	opts := defaultOpts()

	// Apply external options
	for _, f := range optFuncs {
		f(opts)
	}

	c := &http.Client{
		Timeout: opts.Timeout,
	}

	// Instantiate trace provider
	tracer := otel.GetTracerProvider().Tracer(semconv.HttpClientName)

	// Instantiate meter provider
	meter := otel.GetMeterProvider().Meter(semconv.HttpClientName)

	// Instantiate propagator
	propagator := otel.GetTextMapPropagator()

	// Create HTTP client latency histogram
	latency, err := meter.Float64Histogram(
		semconv.HttpClientLatencyName,
		metric.WithUnit("ms"),
		metric.WithDescription("Measures the duration of HTTP request handling"),
		metric.WithExplicitBucketBoundaries(semconv.HttpExplicitBucketBoundaries...),
	)
	if err != nil {
		panic(err)
	}

	return &HttpClient{
		client: c,

		tracer:     tracer,
		meter:      meter,
		propagator: propagator,

		latency: latency,
	}
}

// Configure timeout of the HTTP client
func WithTimeout(timeout time.Duration) OptFunc {
	return func(opts *Opts) {
		opts.Timeout = timeout
	}
}

func (c *HttpClient) Do(
	ctx context.Context,
	req *http.Request,
	spanName string,
) (
	*http.Response,
	error,
) {
	// Start timer
	requestStartTime := time.Now()

	// Create latency recording function
	recordLatencyFunc := func(
		ctx context.Context,
		startTime time.Time,
		attrs []attribute.KeyValue,
	) {
		// Record server latency
		elapsedTime := float64(time.Since(requestStartTime)) / float64(time.Millisecond)
		c.latency.Record(ctx, elapsedTime, metric.WithAttributes(attrs...))
	}

	// Parse HTTP attributes from the request for both span and metrics
	spanAttrs, metricAttrs := c.getSpanAndMetricClientAttributes(req)

	// Create span options
	spanOpts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindClient),
		trace.WithAttributes(spanAttrs...),
	}

	// Start HTTP server span
	ctx, span := c.tracer.Start(ctx, spanName, spanOpts...)
	defer span.End()

	// Inject context into the HTTP headers
	headers := http.Header{}
	carrier := propagation.HeaderCarrier(headers)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	for k, v := range headers {
		req.Header.Add(k, v[0])
	}

	// Perform HTTP request
	res, err := c.client.Do(req)
	if err != nil {
		// Record error
		span.RecordError(err)
		// Record server latency
		recordLatencyFunc(ctx, requestStartTime, metricAttrs)
		return nil, err
	}

	// Add HTTP status code to the attributes
	span.SetAttributes(semconv.HttpResponseStatusCode.Int(res.StatusCode))
	metricAttrs = append(metricAttrs, semconv.HttpResponseStatusCode.Int(res.StatusCode))

	// Record server latency
	recordLatencyFunc(ctx, requestStartTime, metricAttrs)

	return res, err
}

func (m *HttpClient) getSpanAndMetricClientAttributes(
	r *http.Request,
) (
	[]attribute.KeyValue,
	[]attribute.KeyValue,
) {
	spanAttrs := semconv.WithHttpClientAttributes(r)
	metricAttrs := make([]attribute.KeyValue, len(spanAttrs))
	copy(metricAttrs, spanAttrs)

	var url string
	if r.URL != nil {
		// Remove any username/password info that may be in the URL.
		userinfo := r.URL.User
		r.URL.User = nil
		url = r.URL.String()
		// Restore any username/password info that was removed.
		r.URL.User = userinfo
	}
	spanAttrs = append(spanAttrs, semconv.HttpUrlFull.String(url))

	return spanAttrs, metricAttrs
}
