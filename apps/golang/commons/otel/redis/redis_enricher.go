package redis

import (
	"context"
	"strconv"

	semconv "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const redis = "redis"

type redisOpts struct {
	TracerName string
	Server     string
	Port       int
}

type redisOptFunc func(*redisOpts)

type RedisEnricher struct {
	Opts *redisOpts
}

// Create a Redis database instance
func NewRedisEnricher(
	optFuncs ...redisOptFunc,
) *RedisEnricher {

	// Apply external options
	opts := &redisOpts{
		TracerName: "",
		Server:     "",
		Port:       0,
	}
	for _, f := range optFuncs {
		f(opts)
	}

	return &RedisEnricher{
		Opts: opts,
	}
}

// Configure tracer name
func WithTracerName(tracerName string) redisOptFunc {
	return func(opts *redisOpts) {
		opts.TracerName = tracerName
	}
}

// Configure Redis server
func WithServer(server string) redisOptFunc {
	return func(opts *redisOpts) {
		opts.Server = server
	}
}

// Configure Redis port
func WithPort(port string) redisOptFunc {
	return func(opts *redisOpts) {
		p, _ := strconv.Atoi(port)
		opts.Port = p
	}
}

func (e *RedisEnricher) CreateSpan(
	ctx context.Context,
	parentSpan trace.Span,
	operation string,
	key string,
	value string,
) (
	context.Context,
	trace.Span,
) {
	// Create database span
	ctx, dbSpan := parentSpan.TracerProvider().
		Tracer(e.Opts.TracerName).
		Start(
			ctx,
			operation+" "+key,
			trace.WithSpanKind(trace.SpanKindClient),
		)

	// Set additional span attributes
	statement := operation + " " + key + " " + value
	dbSpanAttrs := e.getCommonAttributes()
	dbSpanAttrs = append(dbSpanAttrs, semconv.DatabaseDbStatement.String(statement))
	dbSpan.SetAttributes(dbSpanAttrs...)

	return ctx, dbSpan
}

func (e *RedisEnricher) getCommonAttributes() []attribute.KeyValue {
	return []attribute.KeyValue{
		semconv.ServerAddress.String(e.Opts.Server),
		semconv.ServerPort.Int(e.Opts.Port),
		semconv.DatabaseSystem.String(redis),
	}
}
