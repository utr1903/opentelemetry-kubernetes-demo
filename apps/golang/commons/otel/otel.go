package otel

import (
	"context"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Creates new trace provider
func NewTraceProvider(
	ctx context.Context,
) *sdktrace.TracerProvider {

	// Instantiate OTLP grpc exporter for traces
	exp, err := otlptracegrpc.New(ctx)
	if err != nil {
		panic(err)
	}

	// Instantiate OTel resource
	r := resource.Default()
	if err != nil {
		panic(err)
	}

	// Create trace provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)

	// Set global trace provider
	otel.SetTracerProvider(tp)

	// Set trace propagator
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		))

	return tp
}

// Shuts down trace provider
func ShutdownTraceProvider(
	ctx context.Context,
	tp *sdktrace.TracerProvider,
) {
	// Do not make the application hang when it is shutdown.
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := tp.Shutdown(ctx); err != nil {
		panic(err)
	}
}

// Creates new meter provider
func NewMetricProvider(
	ctx context.Context,
) *sdkmetric.MeterProvider {

	// Instantiate OTLP grpc exporter for metrics
	exp, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		panic(err)
	}

	// Instantiate meter provider
	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp)))

	// Set global meter provider
	otel.SetMeterProvider(mp)
	return mp
}

// Shuts down meter provider
func ShutdownMetricProvider(
	ctx context.Context,
	mp *sdkmetric.MeterProvider,
) {
	// Do not make the application hang when it is shutdown.
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := mp.Shutdown(ctx); err != nil {
		panic(err)
	}
}

// Starts runtime metric collection
func StartCollectingRuntimeMetrics() {
	err := runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	if err != nil {
		panic(err)
	}
}
