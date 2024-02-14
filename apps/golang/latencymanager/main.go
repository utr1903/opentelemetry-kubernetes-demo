package main

import (
	"context"

	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/logger"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/redis"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/latencymanager/config"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/latencymanager/manager"
)

func main() {
	// Create new config
	cfg := config.NewConfig()

	// Initialize logger
	log := logger.NewLogger(cfg.ServiceName)

	// Get context
	ctx := context.Background()

	// Create tracer provider
	tp := otel.NewTraceProvider(ctx)
	defer otel.ShutdownTraceProvider(ctx, tp)

	// Instantiate Redis database
	rdb := redis.New(
		redis.WithServer(cfg.RedisServer),
		redis.WithPort(cfg.RedisPort),
		redis.WithPassword(cfg.RedisPassword),
	)
	rdb.CreateDatabaseConnection()
	defer rdb.Instance.Close()

	// Instantiate and run latency lmgr
	lmgr := manager.NewLatencyManager(log, rdb, cfg)
	lmgr.Run()
}
