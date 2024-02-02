package logger

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

type Logger struct {
	serviceName string
	logger      *logrus.Logger
}

// Creates new logger
func NewLogger(
	serviceName string,
) *Logger {
	return &Logger{
		serviceName: serviceName,
		logger: &logrus.Logger{
			Formatter: new(logrus.JSONFormatter),
			Level:     logrus.InfoLevel,
		},
	}
}

// Logs a message with trace context
func (l *Logger) Log(
	lvl logrus.Level,
	ctx context.Context,
	user string,
	msg string,
) {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().HasTraceID() && span.SpanContext().HasSpanID() {
		l.logger.WithFields(logrus.Fields{
			"user":         user,
			"service.name": l.serviceName,
			"trace.id":     span.SpanContext().TraceID().String(),
			"span.id":      span.SpanContext().SpanID().String(),
		}).Log(lvl, msg)
	} else {
		l.logger.WithFields(logrus.Fields{
			"user": user,
		}).Log(lvl, msg)
	}
}
