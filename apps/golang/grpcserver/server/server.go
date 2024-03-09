package server

import (
	"context"

	pb "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/grpc/proto"
)

type GrpcServer struct {
	pb.GrpcServer
}

func New() *GrpcServer {
	return &GrpcServer{}
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
