package kafkaproducer

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/dtos"
	commonerr "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/error"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/logger"

	otelkafka "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel/kafka"
)

var (
	randomErrors = map[int]string{
		1: commonerr.KAFKA_CONNECTION_ERROR,
		2: commonerr.DATABASE_CONNECTION_ERROR,
		3: commonerr.TABLE_DOES_NOT_EXIST_ERROR,
	}
)

type Opts struct {
	ServiceName     string
	RequestInterval int64
	BrokerAddress   string
	BrokerTopic     string
}

type OptFunc func(*Opts)

func defaultOpts() *Opts {
	return &Opts{
		RequestInterval: 2000,
		BrokerAddress:   "kafka",
		BrokerTopic:     "otel",
	}
}

type KafkaConsumerSimulator struct {
	logger     *logger.Logger
	Opts       *Opts
	Randomizer *rand.Rand
}

// Create an kafka consumer simulator instance
func New(
	log *logger.Logger,
	optFuncs ...OptFunc,
) *KafkaConsumerSimulator {

	// Instantiate options with default values
	opts := defaultOpts()

	// Apply external options
	for _, f := range optFuncs {
		f(opts)
	}

	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))

	return &KafkaConsumerSimulator{
		logger:     log,
		Opts:       opts,
		Randomizer: randomizer,
	}
}

// Configure service name of simulator
func WithServiceName(serviceName string) OptFunc {
	return func(opts *Opts) {
		opts.ServiceName = serviceName
	}
}

// Configure Kafka request interval
func WithRequestInterval(requestInterval string) OptFunc {
	interval, err := strconv.ParseInt(requestInterval, 10, 64)
	if err != nil {
		panic(err.Error())
	}
	return func(opts *Opts) {
		opts.RequestInterval = interval
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

// Starts simulating Kafka consumer
func (k *KafkaConsumerSimulator) Simulate(
	users []string,
) {

	// Create producer
	producer := k.createKafkaProducer()

	// Wrap OTel around the producer
	otelproducer := otelkafka.NewKafkaProducer(producer)

	// Publish messages
	go k.publishMessages(otelproducer, users)
}

// Creates the Kafka producer
func (k *KafkaConsumerSimulator) createKafkaProducer() sarama.AsyncProducer {

	// Create config
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V3_0_0_0
	saramaConfig.Producer.Return.Successes = true
	saramaConfig.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	// Create producer
	producer, err := sarama.NewAsyncProducer(
		[]string{k.Opts.BrokerAddress},
		saramaConfig,
	)
	if err != nil {
		panic(err)
	}

	// Print errors if message publishing goes wrong
	go func() {
		for err := range producer.Errors() {
			fmt.Println("Failed to write message: " + err.Error())
		}
	}()

	return producer
}

// Publish messages to topic
func (k *KafkaConsumerSimulator) publishMessages(
	otelproducer *otelkafka.KafkaProducer,
	users []string,
) {

	// Keep publishing messages
	for {
		func() {
			// Inject tracing info into message
			ctx := context.Background()

			// Make request after each interval
			time.Sleep(time.Duration(k.Opts.RequestInterval) * time.Millisecond)

			// Get a random user
			user := users[k.Randomizer.Intn(len(users))]

			// Create random error
			errType := k.createRandomError()

			// Create message
			body, err := k.createMessageBody(user, errType)
			if err != nil {
				k.logger.Log(logrus.ErrorLevel, ctx, user, "Creating message body failed:"+err.Error())
				return
			}
			msg := sarama.ProducerMessage{
				Topic: k.Opts.BrokerTopic,
				Value: sarama.ByteEncoder(body),
			}

			// Publish message
			k.logger.Log(logrus.InfoLevel, ctx, user, "Publishing message...")
			otelproducer.Publish(ctx, &msg, errType)
			k.logger.Log(logrus.InfoLevel, ctx, user, "Message published successfully.")
		}()
	}
}

func (k *KafkaConsumerSimulator) createRandomError() *string {
	randomNum := k.Randomizer.Intn(15)
	if randomNum == 1 || randomNum == 2 || randomNum == 3 {
		errType := randomErrors[randomNum]
		return &errType
	}
	return nil
}

// Creates the message body with a potential random error
func (k *KafkaConsumerSimulator) createMessageBody(
	user string,
	errType *string,
) (
	[]byte,
	error,
) {

	// Create dto
	dto := &dtos.CreateRequestDto{
		Name: user,
	}

	// Add error if created
	if errType != nil {
		dto.Error = *errType
	}

	return json.Marshal(dto)
}
