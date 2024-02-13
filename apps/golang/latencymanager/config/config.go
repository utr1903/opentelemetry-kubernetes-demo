package config

import (
	"os"
)

type LatencyManagerConfig struct {

	// App name
	ServiceName string

	// Cron job schedule
	CronJobSchedule string

	// Redis
	RedisServer   string
	RedisPort     string
	RedisPassword string
}

// Creates new config object by parsing environment variables
func NewConfig() *LatencyManagerConfig {
	return &LatencyManagerConfig{
		ServiceName: os.Getenv("OTEL_SERVICE_NAME"),

		CronJobSchedule: os.Getenv("CRON_JOB_SCHEDULE"),

		RedisServer:   os.Getenv("REDIS_SERVER"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
	}
}
