package server

import (
	"context"

	"github.com/sirupsen/logrus"
	pb "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/grpc/proto"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/logger"
)

type GrpcServer struct {
	pb.GrpcServer
	logger *logger.Logger
}

func New(
	log *logger.Logger,
) *GrpcServer {
	return &GrpcServer{
		logger: log,
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
	res := &pb.Response{
		IsSucceeded: true,
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
	res := &pb.Response{
		IsSucceeded: true,
	}

	s.logger.Log(logrus.InfoLevel, ctx, "", "Delete method is performed...")
	return res, nil
}
