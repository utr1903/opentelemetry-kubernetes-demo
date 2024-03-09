package grpcclient

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	pb "github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/grpc/proto"
	"github.com/utr1903/opentelemetry-kubernetes-demo/apps/golang/commons/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Opts struct {
	ServiceName     string
	RequestInterval int64
	ServerEndpoint  string
	ServerPort      string
}

type OptFunc func(*Opts)

func defaultOpts() *Opts {
	return &Opts{
		RequestInterval: 2000,
		ServerEndpoint:  "grpcserver",
		ServerPort:      "8080",
	}
}

type GrpcServerSimulator struct {
	logger     *logger.Logger
	Opts       *Opts
	Randomizer *rand.Rand
}

// Create an gRPC server simulator instance
func New(
	log *logger.Logger,
	optFuncs ...OptFunc,
) *GrpcServerSimulator {

	// Instantiate options with default values
	opts := defaultOpts()

	// Apply external options
	for _, f := range optFuncs {
		f(opts)
	}

	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))

	return &GrpcServerSimulator{
		logger:     log,
		Opts:       opts,
		Randomizer: randomizer,
	}
}

// Configure service name of simulator
func WithServiceName(serviceName string) OptFunc {
	return func(opts *Opts) {
		opts.ServiceName = serviceName
	}
}

// Configure HTTP server request interval
func WithRequestInterval(requestInterval string) OptFunc {
	interval, err := strconv.ParseInt(requestInterval, 10, 64)
	if err != nil {
		panic(err.Error())
	}
	return func(opts *Opts) {
		opts.RequestInterval = interval
	}
}

// Configure HTTP server endpoint
func WithServerEndpoint(serverEndpoint string) OptFunc {
	return func(opts *Opts) {
		opts.ServerEndpoint = serverEndpoint
	}
}

// Configure HTTP server port
func WithServerPort(serverPort string) OptFunc {
	return func(opts *Opts) {
		opts.ServerPort = serverPort
	}
}

// Starts simulating gRPC server
func (g *GrpcServerSimulator) Simulate() {

	// Get context
	ctx := context.Background()

	// Create connection to gRPC server
	conn, err := grpc.Dial(
		g.Opts.ServerEndpoint+":"+g.Opts.ServerPort,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		g.logger.Log(logrus.ErrorLevel, ctx, "", "Creating gRPC server connection is failed.")
		panic(err)
	}
	defer conn.Close()

	client := pb.NewGrpcClient(conn)

	// LIST simulator
	go func() {
		for {

			// Make request after each interval
			time.Sleep(time.Duration(g.Opts.RequestInterval) * time.Millisecond)

			_, err := client.Get(ctx, &pb.Request{})
			if err != nil {
				g.logger.Log(logrus.ErrorLevel, ctx, "", "List method is failed.")
			}
		}
	}()

	// DELETE simulator
	go func() {
		for {

			// Make request after each interval
			time.Sleep(time.Duration(g.Opts.RequestInterval) * time.Millisecond)

			_, err := client.Delete(ctx, &pb.Request{})
			if err != nil {
				g.logger.Log(logrus.ErrorLevel, ctx, "", "Delete method is failed.")
			}
		}
	}()
}
