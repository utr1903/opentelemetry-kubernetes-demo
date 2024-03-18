package config

import (
	"os"
)

type GrpcServerConfig struct {

	// App name
	ServiceName string

	// App port
	ServicePort string

	// Redis
	RedisServer   string
	RedisPort     string
	RedisPassword string
}

// Creates new config object by parsing environment variables
func NewConfig() *GrpcServerConfig {
	return &GrpcServerConfig{
		ServiceName: os.Getenv("OTEL_SERVICE_NAME"),
		ServicePort: os.Getenv("APP_PORT"),

		RedisServer:   os.Getenv("REDIS_SERVER"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}
}
