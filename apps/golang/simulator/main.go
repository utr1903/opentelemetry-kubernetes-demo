package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/logger"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/simulator/config"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/simulator/httpclient"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/simulator/kafkaproducer"
)

func main() {
	// Get context
	ctx := context.Background()

	// Create new config
	cfg := config.NewConfig()

	// Initialize logger
	logger.NewLogger(cfg.ServiceName)

	// Create tracer provider
	tp := otel.NewTraceProvider(ctx)
	defer otel.ShutdownTraceProvider(ctx, tp)

	// Create metric provider
	mp := otel.NewMetricProvider(ctx)
	defer otel.ShutdownMetricProvider(ctx, mp)

	// Collect runtime metrics
	otel.StartCollectingRuntimeMetrics()

	// Simulate
	go simulateHttpServer(cfg)
	go simulateKafkaConsumer(cfg)

	// Wait for signal to shutdown the simulator
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	<-ctx.Done()
}

func simulateHttpServer(
	cfg *config.SimulatorConfig,
) {
	// Instantiate HTTP server simulator
	httpserverSimulator := httpclient.New(
		httpclient.WithServiceName(cfg.ServiceName),
		httpclient.WithRequestInterval(cfg.HttpserverRequestInterval),
		httpclient.WithServerEndpoint(cfg.HttpserverEndpoint),
		httpclient.WithServerPort(cfg.HttpserverPort),
	)

	// Simulate
	httpserverSimulator.Simulate(cfg.Users)
}

func simulateKafkaConsumer(
	cfg *config.SimulatorConfig,
) {
	// Instantiate Kafka consumer simulator
	kafkaConsumerSimulator := kafkaproducer.New(
		kafkaproducer.WithServiceName(cfg.ServiceName),
		kafkaproducer.WithRequestInterval(cfg.KafkaRequestInterval),
		kafkaproducer.WithBrokerAddress(cfg.KafkaBrokerAddress),
		kafkaproducer.WithBrokerTopic(cfg.KafkaTopic),
	)

	// Simulate
	kafkaConsumerSimulator.Simulate(cfg.Users)
}
