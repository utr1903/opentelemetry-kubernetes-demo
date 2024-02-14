package config

import (
	"os"
)

type LatencyManagerConfig struct {

	// App name
	ServiceName string
	// Project name
	ClusterName string

	// Cron job schedule
	CronJobSchedule string

	// Redis
	RedisServer   string
	RedisPort     string
	RedisPassword string

	// Observability backend
	ObservabilityBackendName     string
	ObservabilityBackendEndpoint string
	ObservabilityBackendApiKey   string
}

// Creates new config object by parsing environment variables
func NewConfig() *LatencyManagerConfig {
	return &LatencyManagerConfig{
		ServiceName: os.Getenv("OTEL_SERVICE_NAME"),
		ClusterName: os.Getenv("CLUSTER_NAME"),

		CronJobSchedule: os.Getenv("CRON_JOB_SCHEDULE"),

		RedisServer:   os.Getenv("REDIS_SERVER"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),

		ObservabilityBackendName:     os.Getenv("OBSERVABILITY_BACKEND_NAME"),
		ObservabilityBackendEndpoint: os.Getenv("OBSERVABILITY_BACKEND_ENDPOINT"),
		ObservabilityBackendApiKey:   os.Getenv("OBSERVABILITY_BACKEND_API_KEY"),
	}
}
