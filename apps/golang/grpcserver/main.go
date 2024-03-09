package main

import (
	"context"
	"net"

	"github.com/sirupsen/logrus"
	pb "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/grpc/proto"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/logger"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/grpcserver/config"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/grpcserver/server"
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

	// Create server
	lis, err := net.Listen("tcp", ":"+cfg.ServicePort)
	if err != nil {
		log.Log(logrus.ErrorLevel, ctx, "", "Failed to create listener!")
		panic(err)
	}

	// Instantiate gRPC server
	grpcsrv := grpc.NewServer()

	// Instantiate & register server implementation
	srv := server.New(log)
	pb.RegisterGrpcServer(grpcsrv, srv)

	// Start gRPC server
	err = grpcsrv.Serve(lis)
	if err != nil {
		log.Log(logrus.ErrorLevel, ctx, "", "Failed to create gRPC server!")
		panic(err)
	}
}
