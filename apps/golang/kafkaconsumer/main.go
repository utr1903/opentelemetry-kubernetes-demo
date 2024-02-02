package main

import (
	"context"
	"os"
	"os/signal"

	_ "github.com/go-sql-driver/mysql"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/logger"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/mysql"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/kafkaconsumer/config"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/kafkaconsumer/consumer"
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
	db := mysql.New(
		mysql.WithServer(cfg.MysqlServer),
		mysql.WithPort(cfg.MysqlPort),
		mysql.WithUsername(cfg.MysqlUsername),
		mysql.WithPassword(cfg.MysqlPassword),
		mysql.WithDatabase(cfg.MysqlDatabase),
		mysql.WithTable(cfg.MysqlTable),
	)
	db.CreateDatabaseConnection()
	defer db.Instance.Close()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	// Instantiate Kafka consumer
	kafkaConsumer := consumer.New(log, db,
		consumer.WithServiceName(cfg.ServiceName),
		consumer.WithBrokerAddress(cfg.KafkaBrokerAddress),
		consumer.WithBrokerTopic(cfg.KafkaTopic),
		consumer.WithConsumerGroupId(cfg.KafkaGroupId),
	)
	if err := kafkaConsumer.StartConsumerGroup(ctx); err != nil {
		panic(err.Error())
	}

	<-ctx.Done()
}
