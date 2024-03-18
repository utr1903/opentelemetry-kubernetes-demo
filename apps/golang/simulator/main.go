package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/logger"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/simulator/config"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/simulator/grpcclient"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/simulator/httpclient"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/simulator/kafkaproducer"
)

func main() {
	// Get context
	ctx := context.Background()

	// Create new config
	cfg := config.NewConfig()

	// Initialize logger
	log := logger.NewLogger(cfg.ServiceName)

	// Create tracer provider
	tp := otel.NewTraceProvider(ctx)
	defer otel.ShutdownTraceProvider(ctx, tp)

	// Create metric provider
	mp := otel.NewMetricProvider(ctx)
	defer otel.ShutdownMetricProvider(ctx, mp)

	// Collect runtime metrics
	otel.StartCollectingRuntimeMetrics()

	// Simulate
	go simulateKafkaConsumer(cfg, log)
	go simulateHttpServer(cfg, log)
	go simulateGrpcServer(cfg, log)

	// Wait for signal to shutdown the simulator
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	<-ctx.Done()
}

func simulateKafkaConsumer(
	cfg *config.SimulatorConfig,
	log *logger.Logger,
) {
	// Instantiate Kafka consumer simulator
	kafkaConsumerSimulator := kafkaproducer.New(
		log,
		kafkaproducer.WithServiceName(cfg.ServiceName),
		kafkaproducer.WithRequestInterval(cfg.KafkaRequestInterval),
		kafkaproducer.WithBrokerAddress(cfg.KafkaBrokerAddress),
		kafkaproducer.WithBrokerTopic(cfg.KafkaTopic),
	)

	// Simulate
	kafkaConsumerSimulator.Simulate(cfg.Users)
}

func simulateHttpServer(
	cfg *config.SimulatorConfig,
	log *logger.Logger,
) {
	// Instantiate HTTP server simulator
	httpserverSimulator := httpclient.New(
		log,
		httpclient.WithServiceName(cfg.ServiceName),
		httpclient.WithRequestInterval(cfg.HttpserverRequestInterval),
		httpclient.WithServerEndpoint(cfg.HttpserverEndpoint),
		httpclient.WithServerPort(cfg.HttpserverPort),
	)

	// Simulate
	httpserverSimulator.Simulate(cfg.Users)
}

func simulateGrpcServer(
	cfg *config.SimulatorConfig,
	log *logger.Logger,
) {
	// Instantiate HTTP server simulator
	grpcserverSimulator := grpcclient.New(
		log,
		grpcclient.WithServiceName(cfg.ServiceName),
		grpcclient.WithRequestInterval(cfg.GrpcserverRequestInterval),
		grpcclient.WithServerEndpoint(cfg.GrpcserverEndpoint),
		grpcclient.WithServerPort(cfg.GrpcserverPort),
	)

	// Simulate
	grpcserverSimulator.Simulate(cfg.Users)
}
