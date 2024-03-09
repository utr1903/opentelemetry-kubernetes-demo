package config

import "os"

type SimulatorConfig struct {

	// App name
	ServiceName string

	// Kafka producer
	KafkaRequestInterval string
	KafkaBrokerAddress   string
	KafkaTopic           string

	// HTTP server
	HttpserverRequestInterval string
	HttpserverEndpoint        string
	HttpserverPort            string

	// gRPC server
	GrpcserverRequestInterval string
	GrpcserverEndpoint        string
	GrpcserverPort            string

	// Users
	Users []string
}

// Creates new config object by parsing environment variables
func NewConfig() *SimulatorConfig {
	return &SimulatorConfig{
		ServiceName: os.Getenv("OTEL_SERVICE_NAME"),

		KafkaRequestInterval: os.Getenv("KAFKA_REQUEST_INTERVAL"),
		KafkaBrokerAddress:   os.Getenv("KAFKA_BROKER_ADDRESS"),
		KafkaTopic:           os.Getenv("KAFKA_TOPIC"),

		HttpserverRequestInterval: os.Getenv("HTTP_SERVER_REQUEST_INTERVAL"),
		HttpserverEndpoint:        os.Getenv("HTTP_SERVER_ENDPOINT"),
		HttpserverPort:            os.Getenv("HTTP_SERVER_PORT"),

		GrpcserverRequestInterval: os.Getenv("GRPC_SERVER_REQUEST_INTERVAL"),
		GrpcserverEndpoint:        os.Getenv("GRPC_SERVER_ENDPOINT"),
		GrpcserverPort:            os.Getenv("GRPC_SERVER_PORT"),

		Users: []string{
			"elon",
			"jeff",
			"warren",
			"bill",
			"mark",
		},
	}
}
