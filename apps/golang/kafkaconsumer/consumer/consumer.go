package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/dtos"
	commonerr "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/error"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/logger"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/mysql"
	otelkafka "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel/kafka"
	otelmysql "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel/mysql"
	otelredis "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel/redis"
	semconv "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel/semconv/v1.24.0"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/redis"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const CONSUMER string = "kafkaconsumer"

type Opts struct {
	ServiceName     string
	BrokerAddress   string
	BrokerTopic     string
	ConsumerGroupId string
}

type OptFunc func(*Opts)

func defaultOpts() *Opts {
	return &Opts{
		BrokerAddress:   "kafka",
		BrokerTopic:     "otel",
		ConsumerGroupId: "kafkaconsumer",
	}
}

type KafkaConsumer struct {
	logger *logger.Logger
	Opts   *Opts

	MySql             *mysql.MySqlDatabase
	MySqlOtelEnricher *otelmysql.MySqlEnricher

	Redis             *redis.RedisDatabase
	RedisOtelEnricher *otelredis.RedisEnricher
}

// Create a kafka consumer instance
func New(
	log *logger.Logger,
	rdb *redis.RedisDatabase,
	db *mysql.MySqlDatabase,
	optFuncs ...OptFunc,
) *KafkaConsumer {

	// Instantiate options with default values
	opts := defaultOpts()

	// Apply external options
	for _, f := range optFuncs {
		f(opts)
	}

	return &KafkaConsumer{
		logger: log,
		MySql:  db,
		MySqlOtelEnricher: otelmysql.NewMysqlEnricher(
			otelmysql.WithTracerName(CONSUMER),
			otelmysql.WithServer(db.Opts.Server),
			otelmysql.WithPort(db.Opts.Port),
			otelmysql.WithUsername(db.Opts.Username),
			otelmysql.WithDatabase(db.Opts.Database),
			otelmysql.WithTable(db.Opts.Table),
		),
		Redis: rdb,
		RedisOtelEnricher: otelredis.NewRedisEnricher(
			otelredis.WithTracerName(CONSUMER),
			otelredis.WithServer(rdb.Opts.Server),
			otelredis.WithPort(rdb.Opts.Port),
		),
		Opts: opts,
	}
}

// Configure service name of consumer
func WithServiceName(serviceName string) OptFunc {
	return func(opts *Opts) {
		opts.ServiceName = serviceName
	}
}

// Configure Kafka broker address
func WithBrokerAddress(address string) OptFunc {
	return func(opts *Opts) {
		opts.BrokerAddress = address
	}
}

// Configure Kafka broker topic
func WithBrokerTopic(topic string) OptFunc {
	return func(opts *Opts) {
		opts.BrokerTopic = topic
	}
}

// Configure Kafka consumer group ID
func WithConsumerGroupId(groupId string) OptFunc {
	return func(opts *Opts) {
		opts.ConsumerGroupId = groupId
	}
}

func (k *KafkaConsumer) StartConsumerGroup(
	ctx context.Context,
) error {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	saramaConfig.Version = sarama.V3_0_0_0
	saramaConfig.Producer.Return.Successes = true

	consumerGroup, err := sarama.NewConsumerGroup(
		[]string{k.Opts.BrokerAddress},
		k.Opts.ConsumerGroupId,
		saramaConfig,
	)
	if err != nil {
		return err
	}

	otelconsumer := otelkafka.NewKafkaConsumer()
	handler := groupHandler{
		logger:            k.logger,
		ready:             make(chan bool),
		Opts:              k.Opts,
		MySql:             k.MySql,
		MySqlOtelEnricher: k.MySqlOtelEnricher,
		Redis:             k.Redis,
		RedisOtelEnricher: k.RedisOtelEnricher,
		Consumer:          otelconsumer,
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() error {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := consumerGroup.Consume(
				ctx,
				[]string{k.Opts.BrokerTopic},
				&handler,
			); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return err
				}
				fmt.Println("Error from consumer: " + err.Error())
				return err
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				fmt.Println("Context cancelled: " + err.Error())
				return err
			}
			handler.ready = make(chan bool)
		}
	}()

	<-handler.ready // Await till the consumer has been set up

	return nil
}

type groupHandler struct {
	logger            *logger.Logger
	ready             chan bool
	Opts              *Opts
	MySql             *mysql.MySqlDatabase
	MySqlOtelEnricher *otelmysql.MySqlEnricher
	Redis             *redis.RedisDatabase
	RedisOtelEnricher *otelredis.RedisEnricher
	Consumer          *otelkafka.KafkaConsumer
}

func (g *groupHandler) Setup(_ sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(g.ready)
	return nil
}

func (g *groupHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (g *groupHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for {
		select {
		case msg := <-claim.Messages():
			g.consumeMessage(session, msg)

		case <-session.Context().Done():
			return nil
		}
	}
}

func (g *groupHandler) consumeMessage(
	session sarama.ConsumerGroupSession,
	msg *sarama.ConsumerMessage,
) error {

	// Define consume function
	consumeFunc := func(ctx context.Context) error {

		// Parse name out of the message
		body, err := g.parseMessageBody(ctx, msg.Value)
		if err != nil {
			g.logger.Log(logrus.ErrorLevel, ctx, "", "Consuming message is failed: "+err.Error())
			return err
		}

		name := body.Name
		errType := body.Error

		g.logger.Log(logrus.InfoLevel, ctx, name, "Consuming message...")

		// Put it into Redis
		err = g.setIntoRedis(ctx, name)
		if err != nil {
			g.logger.Log(logrus.ErrorLevel, ctx, name, "Consuming message is failed: "+err.Error())
			return err
		}

		// Store it into MySQL
		err = g.storeIntoMysql(ctx, name, errType)
		if err != nil {
			g.logger.Log(logrus.ErrorLevel, ctx, name, "Consuming message is failed: "+err.Error())
			return err
		}

		// Acknowledge message
		session.MarkMessage(msg, "")
		g.logger.Log(logrus.InfoLevel, ctx, name, "Consuming message is succeeded.")

		return nil
	}

	// Execute consume within OTel wrapper
	ctx := context.Background()
	g.Consumer.Consume(ctx, msg, g.Opts.ConsumerGroupId, consumeFunc)

	return nil
}

func (g *groupHandler) parseMessageBody(
	ctx context.Context,
	messageBody []byte,
) (
	*dtos.CreateRequestDto,
	error,
) {

	// Start parsing span
	parentSpan := trace.SpanFromContext(ctx)
	ctx, parseSpan := parentSpan.TracerProvider().
		Tracer(semconv.KafkaConsumerName).
		Start(
			ctx,
			"parse dto",
			trace.WithSpanKind(trace.SpanKindInternal),
		)
	defer parseSpan.End()

	g.logger.Log(logrus.InfoLevel, ctx, "", "Parsing dto...")

	dto := &dtos.CreateRequestDto{}
	err := json.Unmarshal(messageBody, dto)
	if err != nil {
		msg := "Parsing dto failed."
		g.logger.Log(logrus.ErrorLevel, ctx, "", msg)
		g.addErrorToSpan(parseSpan, msg, err)
		return nil, err
	}

	g.logger.Log(logrus.InfoLevel, ctx, dto.Name, "Parsing dto succeeded.")
	return dto, nil
}

func (g *groupHandler) setIntoRedis(
	ctx context.Context,
	name string,
) error {

	// Create database span
	parentSpan := trace.SpanFromContext(ctx)
	_, dbSpan := g.RedisOtelEnricher.CreateSpan(
		ctx,
		parentSpan,
		"SET",
		"name",
	)
	defer dbSpan.End()

	// Set the new latency status
	err := g.Redis.Instance.Set("name", name, time.Hour).Err()
	if err != nil {
		g.logger.Log(logrus.ErrorLevel, ctx, CONSUMER, "Redis variable [name] could not be set: "+err.Error())
		return err
	}
	g.logger.Log(logrus.InfoLevel, ctx, CONSUMER, "Redis variable [name] is set successfully.")
	return nil
}

func (g *groupHandler) storeIntoMysql(
	ctx context.Context,
	name string,
	errType string,
) error {

	g.logger.Log(logrus.InfoLevel, ctx, name, "Storing into DB...")

	// Create table does not exist error
	var dbStatement string
	dbOperation := "INSERT"
	if errType == commonerr.TABLE_DOES_NOT_EXIST_ERROR {
		dbStatement = dbOperation + " INTO " + "faketable" + " (name) VALUES (?)"
	} else {
		dbStatement = dbOperation + " INTO " + g.MySql.Opts.Table + " (name) VALUES (?)"
	}

	// Create database span
	parentSpan := trace.SpanFromContext(ctx)
	ctx, dbSpan := g.MySqlOtelEnricher.CreateSpan(
		ctx,
		parentSpan,
		dbOperation,
		dbStatement,
	)
	defer dbSpan.End()

	// Prepare a statement
	stmt, err := g.MySql.Instance.Prepare(dbStatement)
	if err != nil {
		msg := "Preparing DB statement is failed: " + err.Error()
		g.logger.Log(logrus.ErrorLevel, ctx, name, msg)

		// Add error to span
		g.addErrorToSpan(dbSpan, msg, err)

		return err
	}
	defer stmt.Close()

	// Execute the statement
	_, err = stmt.Exec(name)
	if err != nil {
		msg := "Storing into DB is failed: " + err.Error()
		g.logger.Log(logrus.ErrorLevel, ctx, name, msg)

		// Add error to span
		g.addErrorToSpan(dbSpan, msg, err)

		return err
	}

	// Create database connection error
	if errType == commonerr.DATABASE_CONNECTION_ERROR {
		msg := "Connection to database is lost."
		g.logger.Log(logrus.ErrorLevel, ctx, name, msg)

		// Add error to span
		g.addErrorToSpan(dbSpan, msg, err)

		return errors.New("database connection lost")
	}

	g.logger.Log(logrus.InfoLevel, ctx, name, "Storing into DB is succeeded.")
	return nil
}

// Add error to span
func (g *groupHandler) addErrorToSpan(
	span trace.Span,
	description string,
	err error,
) {

	dbSpanAttrs := []attribute.KeyValue{
		semconv.OtelStatusCode.String("ERROR"),
		semconv.OtelStatusDescription.String(description),
	}
	span.SetAttributes(dbSpanAttrs...)
	span.RecordError(
		err,
		trace.WithAttributes(
			semconv.ExceptionEscaped.Bool(true),
		))
}
