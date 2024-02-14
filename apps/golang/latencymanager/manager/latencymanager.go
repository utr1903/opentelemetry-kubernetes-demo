package manager

import (
	"context"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	commonerr "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/error"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/logger"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/monitoring"
	otelredis "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel/redis"
	semconv "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel/semconv/v1.24.0"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/redis"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/latencymanager/config"
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

	clusterName                  string
	observabilityBackendName     string
	observabilityBackendEndpoint string
	observabilityBackendApiKey   string
}

// Create a HTTP server instance
func NewLatencyManager(
	log *logger.Logger,
	rdb *redis.RedisDatabase,
	cfg *config.LatencyManagerConfig,
) *LatencyManager {

	// Instantiate trace provider
	tracer := otel.GetTracerProvider().Tracer(LATENCY_MANAGER)

	// Instantiate propagator
	propagator := otel.GetTextMapPropagator()

	return &LatencyManager{
		logger:                       log,
		tracer:                       tracer,
		propagator:                   propagator,
		cronJobName:                  cfg.ServiceName,
		cronJobSchedule:              cfg.CronJobSchedule,
		clusterName:                  cfg.ClusterName,
		observabilityBackendName:     cfg.ObservabilityBackendName,
		observabilityBackendEndpoint: cfg.ObservabilityBackendEndpoint,
		observabilityBackendApiKey:   cfg.ObservabilityBackendApiKey,
		Redis:                        rdb,
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

	// Deploy change marker to defined observability backend
	err = m.deployMarker(ctx, isLatencyIncreaseEnabled)
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
	var enabled bool
	increaseLatency, _ := m.Redis.Instance.Get(commonerr.INCREASE_HTTPSERVER_LATENCY).Result()
	if increaseLatency == "true" {
		m.logger.Log(logrus.InfoLevel, ctx, LATENCY_MANAGER, "Redis variable ["+commonerr.INCREASE_HTTPSERVER_LATENCY+"] is found.")
		enabled = true
	} else {
		m.logger.Log(logrus.InfoLevel, ctx, LATENCY_MANAGER, "Redis variable ["+commonerr.INCREASE_HTTPSERVER_LATENCY+"] is not found.")
		enabled = false
	}
	// Create attributes array
	attrs := make([]attribute.KeyValue, 0, 1)
	attrs = append(attrs, attribute.Key("increase.httpserver.latency").String(strconv.FormatBool(enabled)))
	dbSpan.SetAttributes(attrs...)

	return enabled
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
	enableLatencyIncrease := strconv.FormatBool(!isLatencyIncreaseEnabled)
	attrs = append(attrs, attribute.Key("increase.httpserver.latency").String(enableLatencyIncrease))
	dbSpan.SetAttributes(attrs...)

	// Set the new latency status
	err := m.Redis.Instance.Set(commonerr.INCREASE_HTTPSERVER_LATENCY, enableLatencyIncrease, time.Hour).Err()
	if err != nil {
		m.logger.Log(logrus.ErrorLevel, ctx, LATENCY_MANAGER, "Redis variable ["+commonerr.INCREASE_HTTPSERVER_LATENCY+"] could not be set: "+err.Error())
		return err
	}
	m.logger.Log(logrus.InfoLevel, ctx, LATENCY_MANAGER, "Redis variable ["+commonerr.INCREASE_HTTPSERVER_LATENCY+"] is set successfully.")
	return nil
}

func (m *LatencyManager) deployMarker(
	ctx context.Context,
	isLatencyIncreaseEnabled bool,
) error {
	marker := monitoring.NewMarker(m.logger, m.observabilityBackendName, m.observabilityBackendEndpoint, m.observabilityBackendApiKey)
	if marker == nil {
		m.logger.Log(logrus.InfoLevel, ctx, LATENCY_MANAGER, "No observability backend is found for marking changes.")
		return nil
	}

	if isLatencyIncreaseEnabled {
		return marker.Run(
			ctx,
			"httpserver-golang",
			"Rolledback to stable version.",
			"Rollback",
			"Junior developers should not commit to main.",
			m.clusterName,
			"v0.5.3",
		)
	} else {
		return marker.Run(
			ctx,
			"httpserver-golang",
			"Only noobs document changes...",
			"Add mega feature",
			"Life changing feature is added!",
			m.clusterName,
			"v0.6.0",
		)
	}
}
