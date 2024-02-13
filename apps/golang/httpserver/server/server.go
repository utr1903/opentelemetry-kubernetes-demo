package server

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"net/http"

	"github.com/sirupsen/logrus"
	commonerr "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/error"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/logger"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/mysql"
	otelmysql "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel/mysql"
	otelredis "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel/redis"
	semconv "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel/semconv/v1.24.0"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/redis"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const SERVER string = "httpserver"

type Server struct {
	logger *logger.Logger

	MySql             *mysql.MySqlDatabase
	MySqlOtelEnricher *otelmysql.MySqlEnricher

	Redis             *redis.RedisDatabase
	RedisOtelEnricher *otelredis.RedisEnricher
}

// Create a HTTP server instance
func New(
	log *logger.Logger,
	mdb *mysql.MySqlDatabase,
	rdb *redis.RedisDatabase,
) *Server {

	return &Server{
		logger: log,
		MySql:  mdb,
		MySqlOtelEnricher: otelmysql.NewMysqlEnricher(
			otelmysql.WithTracerName(SERVER),
			otelmysql.WithServer(mdb.Opts.Server),
			otelmysql.WithPort(mdb.Opts.Port),
			otelmysql.WithUsername(mdb.Opts.Username),
			otelmysql.WithDatabase(mdb.Opts.Database),
			otelmysql.WithTable(mdb.Opts.Table),
		),
		Redis: rdb,
		RedisOtelEnricher: otelredis.NewRedisEnricher(
			otelredis.WithTracerName(SERVER),
			otelredis.WithServer(rdb.Opts.Server),
			otelredis.WithPort(rdb.Opts.Port),
		),
	}
}

// Liveness
func (s *Server) Livez(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Readiness
func (s *Server) Readyz(
	w http.ResponseWriter,
	r *http.Request,
) {

	// MySQL
	err := s.MySql.Instance.Ping()
	if err != nil {
		s.logger.Log(logrus.ErrorLevel, r.Context(), "MySQL is not reachable.", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Not OK"))
	}

	// Redis
	_, err = s.Redis.Instance.Ping().Result()
	if err != nil {
		s.logger.Log(logrus.ErrorLevel, r.Context(), "Redis is not reachable.", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Not OK"))
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// Server handler
func (s *Server) ServerHandler(
	w http.ResponseWriter,
	r *http.Request,
) {

	// Get server span
	parentSpan := trace.SpanFromContext(r.Context())
	defer parentSpan.End()

	s.logger.Log(logrus.InfoLevel, r.Context(), s.getUser(r), "Server handler is triggered")

	// Perform Redis database query
	s.performRedisQuery(r, parentSpan)

	// Perform MySQL database query
	err := s.performMysqlQuery(w, r, parentSpan)
	if err != nil {
		return
	}

	s.performPostprocessing(r, parentSpan)
	s.createHttpResponse(&w, http.StatusOK, []byte("Success"), parentSpan)
}

// Performs the database query against the Redis database
func (s *Server) performRedisQuery(
	r *http.Request,
	parentSpan trace.Span,
) {
	s.logger.Log(logrus.InfoLevel, r.Context(), s.getUser(r), "Querying Redis...")

	// Create database span
	_, dbSpan := s.RedisOtelEnricher.CreateSpan(
		r.Context(),
		parentSpan,
		"GET",
		commonerr.INCREASE_HTTPSERVER_LATENCY,
	)
	defer dbSpan.End()

	// Retrieve variables from Redis
	increaseLatency, _ := s.Redis.Instance.Get(commonerr.INCREASE_HTTPSERVER_LATENCY).Result()
	if increaseLatency == "true" {
		s.logger.Log(logrus.WarnLevel, r.Context(), s.getUser(r), "Redis variable ["+commonerr.INCREASE_HTTPSERVER_LATENCY+"] is found.")
		time.Sleep(time.Second)
	}
	// Create attributes array
	attrs := make([]attribute.KeyValue, 0, 1)
	attrs = append(attrs, attribute.Key("increase.httpserver.latency").String(increaseLatency))
	dbSpan.SetAttributes(attrs...)
}

// Performs the database query against the MySQL database
func (s *Server) performMysqlQuery(
	w http.ResponseWriter,
	r *http.Request,
	parentSpan trace.Span,
) error {

	// Build query
	dbOperation, dbStatement, err := s.createDbQuery(r)
	if err != nil {
		s.createHttpResponse(&w, http.StatusMethodNotAllowed, []byte("Method not allowed"), parentSpan)
		return err
	}

	// Create database span
	ctx, dbSpan := s.MySqlOtelEnricher.CreateSpan(
		r.Context(),
		parentSpan,
		dbOperation,
		dbStatement,
	)
	defer dbSpan.End()

	// Perform query
	err = s.executeDbQuery(ctx, r, dbStatement)
	if err != nil {
		msg := "Executing DB query is failed."
		s.logger.Log(logrus.ErrorLevel, ctx, s.getUser(r), msg)

		// Add error to span
		s.addErrorToSpan(dbSpan, msg, err)

		s.createHttpResponse(&w, http.StatusInternalServerError, []byte(err.Error()), parentSpan)
		return err
	}

	// Create database connection error
	databaseConnectionError := r.URL.Query().Get(commonerr.DATABASE_CONNECTION_ERROR)
	if databaseConnectionError == "true" {
		msg := "Connection to database is lost."
		s.logger.Log(logrus.ErrorLevel, ctx, s.getUser(r), msg)

		// Add error to span
		s.addErrorToSpan(dbSpan, msg, err)

		s.createHttpResponse(&w, http.StatusInternalServerError, []byte(msg), parentSpan)
		return errors.New("database connection lost")
	}

	return nil
}

// Creates the database query operation and statement
func (s *Server) createDbQuery(
	r *http.Request,
) (
	string,
	string,
	error,
) {
	s.logger.Log(logrus.InfoLevel, r.Context(), s.getUser(r), "Building MySQL query...")

	var dbOperation string
	var dbStatement string

	switch r.Method {
	case http.MethodGet:
		dbOperation = "SELECT"

		// Create table does not exist error
		tableDoesNotExistError := r.URL.Query().Get(commonerr.TABLE_DOES_NOT_EXIST_ERROR)
		if tableDoesNotExistError == "true" {
			dbStatement = dbOperation + " name FROM " + "faketable"
		} else {
			dbStatement = dbOperation + " name FROM " + s.MySql.Opts.Table
		}
		return dbOperation, dbStatement, nil
	case http.MethodDelete:
		dbOperation = "DELETE"
		dbStatement = dbOperation + " FROM " + s.MySql.Opts.Table
	default:
		s.logger.Log(logrus.ErrorLevel, r.Context(), s.getUser(r), "Method is not allowed.")
		return "", "", errors.New("method not allowed")
	}

	s.logger.Log(logrus.InfoLevel, r.Context(), s.getUser(r), "MySQL query is built.")
	return dbOperation, dbStatement, nil
}

// Executes the MySQL database statement
func (s *Server) executeDbQuery(
	ctx context.Context,
	r *http.Request,
	dbStatement string,
) error {

	user := s.getUser(r)
	s.logger.Log(logrus.InfoLevel, ctx, user, "Executing MySQL query...")

	switch r.Method {
	case http.MethodGet:
		// Perform a query
		rows, err := s.MySql.Instance.Query(dbStatement)
		if err != nil {
			s.logger.Log(logrus.ErrorLevel, ctx, user, err.Error())
			return err
		}
		defer rows.Close()

		// Iterate over the results
		names := make([]string, 0, 10)
		for rows.Next() {
			var name string
			err = rows.Scan(&name)
			if err != nil {
				s.logger.Log(logrus.ErrorLevel, ctx, user, err.Error())
				return err
			}
			names = append(names, name)
		}

		_, err = json.Marshal(names)
		if err != nil {
			s.logger.Log(logrus.ErrorLevel, ctx, user, err.Error())
			return err
		}
	case http.MethodDelete:
		_, err := s.MySql.Instance.Exec(dbStatement)
		if err != nil {
			s.logger.Log(logrus.ErrorLevel, ctx, user, err.Error())
			return err
		}
	default:
		s.logger.Log(logrus.ErrorLevel, ctx, user, "Method is not allowed.")
		return errors.New("method not allowed")
	}

	s.logger.Log(logrus.InfoLevel, ctx, user, "MySQL query is executed.")
	return nil
}

// Creates a HTTP response
func (s *Server) createHttpResponse(
	w *http.ResponseWriter,
	statusCode int,
	body []byte,
	serverSpan trace.Span,
) {
	(*w).WriteHeader(statusCode)
	(*w).Write(body)

	attrs := []attribute.KeyValue{
		semconv.HttpResponseStatusCode.Int(statusCode),
	}
	serverSpan.SetAttributes(attrs...)
}

// Performs a postprocessing step
func (s *Server) performPostprocessing(
	r *http.Request,
	parentSpan trace.Span,
) {
	ctx, processingSpan := parentSpan.TracerProvider().
		Tracer(SERVER).
		Start(
			r.Context(),
			"postprocessing",
			trace.WithSpanKind(trace.SpanKindInternal),
		)
	defer processingSpan.End()

	s.produceSchemaNotFoundInCacheWarning(ctx, r)
}

func (s *Server) produceSchemaNotFoundInCacheWarning(
	ctx context.Context,
	r *http.Request,
) {
	s.logger.Log(logrus.InfoLevel, ctx, s.getUser(r), "Postprocessing...")
	schemaNotFoundInCacheWarning := r.URL.Query().Get(commonerr.SCHEMA_NOT_FOUND_IN_CACHE)
	if schemaNotFoundInCacheWarning == "true" {
		user := s.getUser(r)
		s.logger.Log(logrus.WarnLevel, ctx, user, "Processing schema not found in cache. Calculating from scratch.")
		time.Sleep(time.Millisecond * 500)
	} else {
		time.Sleep(time.Millisecond * 10)
	}
	s.logger.Log(logrus.InfoLevel, r.Context(), s.getUser(r), "Postprocessing is complete.")
}

func (s *Server) getUser(
	r *http.Request,
) string {

	user := r.Header.Get("X-User-ID")
	if user == "" {
		user = "_anonymous_"
	}
	return user
}

// Add error to span
func (s *Server) addErrorToSpan(
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
