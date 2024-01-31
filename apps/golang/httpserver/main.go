package main

import (
	"context"

	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/logger"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/httpserver/config"
)

func main() {

	// Create new config
	config.NewConfig()
	cfg := config.GetConfig()

	// Initialize logger
	logger.NewLogger(cfg.ServiceName)

	// Get context
	ctx := context.Background()

	// Create tracer provider
	tp := otel.NewTraceProvider(ctx)
	defer otel.ShutdownTraceProvider(ctx, tp)

	// Create metric provider
	mp := otel.NewMetricProvider(ctx)
	defer otel.ShutdownMetricProvider(ctx, mp)

	// Collect runtime metrics
	otel.StartCollectingRuntimeMetrics()
}
