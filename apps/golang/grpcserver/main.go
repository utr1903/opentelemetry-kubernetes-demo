package main

import (
	"context"
	"net"

	"github.com/sirupsen/logrus"
	pb "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/grpc/proto"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/logger"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/redis"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/grpcserver/config"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/grpcserver/server"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
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

	// Instantiate Redis database
	rdb := redis.New(
		redis.WithServer(cfg.RedisServer),
		redis.WithPort(cfg.RedisPort),
		redis.WithPassword(cfg.RedisPassword),
	)
	rdb.CreateDatabaseConnection()
	defer rdb.Instance.Close()

	// Create server
	lis, err := net.Listen("tcp", ":"+cfg.ServicePort)
	if err != nil {
		log.Log(logrus.ErrorLevel, ctx, "", "Failed to create listener!")
		panic(err)
	}

	log.Log(logrus.InfoLevel, ctx, "", "Listener is created.")

	// Instantiate gRPC server
	grpcsrv := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	log.Log(logrus.InfoLevel, ctx, "", "gRPC server is created.")

	// Instantiate & register server implementation
	srv := server.New(log, rdb)
	pb.RegisterGrpcServer(grpcsrv, srv)

	// Start gRPC server
	err = grpcsrv.Serve(lis)
	if err != nil {
		log.Log(logrus.ErrorLevel, ctx, "", "Failed to run gRPC server!")
		panic(err)
	}
}
