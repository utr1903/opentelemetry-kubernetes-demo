package server

import (
	"context"

	"github.com/sirupsen/logrus"
	pb "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/grpc/proto"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/logger"
	otelredis "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/otel/redis"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/redis"
	"go.opentelemetry.io/otel/trace"
)

const SERVER string = "grpcserver"

type GrpcServer struct {
	pb.GrpcServer

	logger *logger.Logger

	Redis             *redis.RedisDatabase
	RedisOtelEnricher *otelredis.RedisEnricher
}

func New(
	log *logger.Logger,
	rdb *redis.RedisDatabase,
) *GrpcServer {
	return &GrpcServer{
		logger: log,
		Redis:  rdb,
		RedisOtelEnricher: otelredis.NewRedisEnricher(
			otelredis.WithTracerName(SERVER),
			otelredis.WithServer(rdb.Opts.Server),
			otelredis.WithPort(rdb.Opts.Port),
		),
	}
}

func (s *GrpcServer) Get(
	ctx context.Context,
	in *pb.Request,
) (
	*pb.Response,
	error,
) {
	s.logger.Log(logrus.InfoLevel, ctx, "", "List method is triggered...")

	isSucceeded := true
	err := s.performRedisQuery(ctx, false)
	if err != nil {
		isSucceeded = false
	}

	res := &pb.Response{
		IsSucceeded: isSucceeded,
	}

	s.logger.Log(logrus.InfoLevel, ctx, "", "List method is performed...")
	return res, nil
}

func (s *GrpcServer) Delete(
	ctx context.Context,
	in *pb.Request,
) (
	*pb.Response,
	error,
) {
	s.logger.Log(logrus.InfoLevel, ctx, "", "Delete method is triggered...")

	isSucceeded := true
	err := s.performRedisQuery(ctx, true)
	if err != nil {
		isSucceeded = false
	}

	res := &pb.Response{
		IsSucceeded: isSucceeded,
	}

	s.logger.Log(logrus.InfoLevel, ctx, "", "Delete method is performed...")
	return res, nil
}

// Performs the database query against the Redis database
func (s *GrpcServer) performRedisQuery(
	ctx context.Context,
	isDelete bool,
) error {
	s.logger.Log(logrus.InfoLevel, ctx, "", "Querying Redis...")
	parentSpan := trace.SpanFromContext(ctx)

	if !isDelete {
		// Create database span
		ctx, dbSpan := s.RedisOtelEnricher.CreateSpan(
			ctx,
			parentSpan,
			"GET",
			"name",
		)
		defer dbSpan.End()

		// Get name from Redis
		name, err := s.Redis.Instance.Get("name").Result()
		if err != nil {
			s.logger.Log(logrus.ErrorLevel, ctx, "", "Redis variable [name] could not be returned: "+err.Error())
			return err
		}
		s.logger.Log(logrus.InfoLevel, ctx, "", "Redis variable [name] is: "+name)
	} else {
		// Create database span
		ctx, dbSpan := s.RedisOtelEnricher.CreateSpan(
			ctx,
			parentSpan,
			"DEL",
			"name",
		)
		defer dbSpan.End()

		// Delete name from Redis
		_, err := s.Redis.Instance.Del("name").Result()
		if err != nil {
			s.logger.Log(logrus.ErrorLevel, ctx, "", "Redis variable [name] could not be deleted: "+err.Error())
			return err
		}
		s.logger.Log(logrus.InfoLevel, ctx, "", "Redis variable [name] is deleted.")

	}

	return nil
}
