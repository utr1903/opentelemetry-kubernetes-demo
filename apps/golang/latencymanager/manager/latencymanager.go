package manager

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	commonerr "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/error"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/logger"
	otelredis "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel/redis"
	semconv "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel/semconv/v1.24.0"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/redis"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

const LATENCY_MANAGER string = "latencymanager"

type LatencyManager struct {
	logger *logger.Logger

	tracer     trace.Tracer
	propagator propagation.TextMapPropagator

	cronJobName     string
	cronJobSchedule string

	Redis             *redis.RedisDatabase
	RedisOtelEnricher *otelredis.RedisEnricher
}

// Create a HTTP server instance
func NewLatencyManager(
	log *logger.Logger,
	rdb *redis.RedisDatabase,
	cronJobName string,
	cronJobSchedule string,
) *LatencyManager {

	// Instantiate trace provider
	tracer := otel.GetTracerProvider().Tracer(LATENCY_MANAGER)

	// Instantiate propagator
	propagator := otel.GetTextMapPropagator()

	return &LatencyManager{
		logger:          log,
		tracer:          tracer,
		propagator:      propagator,
		cronJobName:     cronJobName,
		cronJobSchedule: cronJobSchedule,
		Redis:           rdb,
		RedisOtelEnricher: otelredis.NewRedisEnricher(
			otelredis.WithTracerName(LATENCY_MANAGER),
			otelredis.WithServer(rdb.Opts.Server),
			otelredis.WithPort(rdb.Opts.Port),
		),
	}
}

func (m *LatencyManager) Run() {

	ctx := context.Background()

	// Start cron job cronJobSpan
	ctx, cronJobSpan := m.tracer.Start(ctx, m.cronJobName, m.getCronJobAttributes()...)
	defer cronJobSpan.End()

	m.logger.Log(logrus.InfoLevel, ctx, LATENCY_MANAGER, "Cron job ["+m.cronJobName+"] is started.")

	// Get current latency value from Redis
	isLatencyIncreaseEnabled := m.getCurrentLatencyStatus(ctx, cronJobSpan)

	// Set new latency value to Redis
	err := m.setNewLatencyStatus(ctx, cronJobSpan, isLatencyIncreaseEnabled)
	if err != nil {
		return
	}

	m.logger.Log(logrus.InfoLevel, ctx, LATENCY_MANAGER, "Cron job ["+m.cronJobName+"] is finished successfully.")
}

func (m *LatencyManager) getCronJobAttributes() []trace.SpanStartOption {

	// Create attributes array
	attrs := make([]attribute.KeyValue, 0, 1)

	// Add attributes
	attrs = append(attrs, semconv.CronJobSchedule.String(m.cronJobSchedule))

	// Create span options
	spanOpts := []trace.SpanStartOption{
		trace.WithSpanKind(trace.SpanKindInternal),
		trace.WithAttributes(attrs...),
	}
	return spanOpts
}

func (m *LatencyManager) getCurrentLatencyStatus(
	ctx context.Context,
	parentSpan trace.Span,
) bool {

	// Create database span
	_, dbSpan := m.RedisOtelEnricher.CreateSpan(
		ctx,
		parentSpan,
		"GET",
		commonerr.INCREASE_HTTPSERVER_LATENCY,
	)
	defer dbSpan.End()

	// Retrieve variables from Redis
	increaseLatency, _ := m.Redis.Instance.Get(commonerr.INCREASE_HTTPSERVER_LATENCY).Result()
	if increaseLatency == "true" {
		m.logger.Log(logrus.InfoLevel, ctx, LATENCY_MANAGER, "Redis variable ["+commonerr.INCREASE_HTTPSERVER_LATENCY+"] is found.")
		return true
	} else {
		m.logger.Log(logrus.InfoLevel, ctx, LATENCY_MANAGER, "Redis variable ["+commonerr.INCREASE_HTTPSERVER_LATENCY+"] is not found.")
		return false
	}
}

func (m *LatencyManager) setNewLatencyStatus(
	ctx context.Context,
	parentSpan trace.Span,
	isLatencyIncreaseEnabled bool,
) error {

	// Create database span
	_, dbSpan := m.RedisOtelEnricher.CreateSpan(
		ctx,
		parentSpan,
		"SET",
		commonerr.INCREASE_HTTPSERVER_LATENCY,
	)
	defer dbSpan.End()

	// Create attributes array
	attrs := make([]attribute.KeyValue, 0, 1)

	// If latency increase is enabled, disable it & vice versa
	var enable string
	if isLatencyIncreaseEnabled {
		enable = "" // disable
	} else {
		enable = "true" // enable
	}

	attrs = append(attrs, attribute.Key("increase.httpserver.latency").String(enable))
	dbSpan.SetAttributes(attrs...)

	// Set the new latency status
	_, err := m.Redis.Instance.Set(commonerr.INCREASE_HTTPSERVER_LATENCY, enable, time.Hour).Result()
	if err != nil {
		m.logger.Log(logrus.ErrorLevel, ctx, LATENCY_MANAGER, "Redis variable ["+commonerr.INCREASE_HTTPSERVER_LATENCY+"] could not be set: "+err.Error())
		return err
	}
	m.logger.Log(logrus.InfoLevel, ctx, LATENCY_MANAGER, "Redis variable ["+commonerr.INCREASE_HTTPSERVER_LATENCY+"] is set successfully.")
	return nil
}
