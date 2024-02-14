package main

import (
	"context"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/logger"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/mysql"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel"
	otelhttp "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel/http"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/redis"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/httpserver/config"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/httpserver/server"
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

	// Create metric provider
	mp := otel.NewMetricProvider(ctx)
	defer otel.ShutdownMetricProvider(ctx, mp)

	// Collect runtime metrics
	otel.StartCollectingRuntimeMetrics()

	// Instantiate MySQL database
	mdb := mysql.New(
		mysql.WithServer(cfg.MysqlServer),
		mysql.WithPort(cfg.MysqlPort),
		mysql.WithUsername(cfg.MysqlUsername),
		mysql.WithPassword(cfg.MysqlPassword),
		mysql.WithDatabase(cfg.MysqlDatabase),
		mysql.WithTable(cfg.MysqlTable),
	)
	mdb.CreateDatabaseConnection()
	defer mdb.Instance.Close()

	// Instantiate Redis database
	rdb := redis.New(
		redis.WithServer(cfg.RedisServer),
		redis.WithPort(cfg.RedisPort),
		redis.WithPassword(cfg.RedisPassword),
	)
	rdb.CreateDatabaseConnection()
	defer rdb.Instance.Close()

	// Instantiate server
	server := server.New(log, mdb, rdb)

	// Serve
	http.Handle("/api", otelhttp.NewHandler(http.HandlerFunc(server.ServerHandler), "api"))
	http.Handle("/livez", http.HandlerFunc(server.Livez))
	http.Handle("/readyz", http.HandlerFunc(server.Readyz))
	http.ListenAndServe(":"+cfg.ServicePort, nil)
}
