package server

import (
	"context"

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
	return &pb.Response{
		IsSucceeded: true,
	}, nil
}

func (s *GrpcServer) Delete(
	ctx context.Context,
	in *pb.Request,
) (
	*pb.Response,
	error,
) {
	return &pb.Response{
		IsSucceeded: true,
	}, nil
}
