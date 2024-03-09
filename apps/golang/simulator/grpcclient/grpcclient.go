package grpcclient

import (
	"context"
	"fmt"

	pb "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcServerSimulator struct {
}

// Create an gRPC server simulator instance
func New() *GrpcServerSimulator {
	return &GrpcServerSimulator{}
}

// Starts simulating gRPC server
func (h *GrpcServerSimulator) Simulate() {
	ctx := context.Background()

	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewGrpcClient(conn)
	res, err := client.Get(ctx, &pb.Request{IsDelete: false})
	if err != nil {
		panic(err)
	}

	fmt.Println(res.IsSucceeded)
}
