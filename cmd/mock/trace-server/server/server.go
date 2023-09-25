package server

import (
	"context"
	"fmt"
	"log"
	"net"

	// status "google.golang.org/grpc/status"
	// codes "google.golang.org/grpc/codes"
	colTrace "go.opentelemetry.io/proto/otlp/collector/trace/v1"

	"os/signal"
	"syscall"

	"github.com/LukaszSwolkien/ingest-tools/cmd/mock/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	colTrace.UnimplementedTraceServiceServer
	server.ServerCore
	ctx          context.Context
	grpcListener net.Listener
	grpcServer   *grpc.Server
	// httpListener	net.Listener
}

func (s *Server) Export(ctx context.Context, in *colTrace.ExportTraceServiceRequest) (*colTrace.ExportTraceServiceResponse, error) {
	log.Printf("Received: %v", in.String())

	// return nil, status.Errorf(codes.Unimplemented, "method Export not implemented")
	return &colTrace.ExportTraceServiceResponse{PartialSuccess: nil}, nil
}

// New creates a new server object
func New(conf server.Conf) *Server {
	s := &Server{
		ctx:        context.Background(),
		grpcServer: grpc.NewServer(),
	}
	s.ServerCore.Init(conf)

	colTrace.RegisterTraceServiceServer(s.grpcServer, &Server{})

	// Register reflection service on gRPC server.
	reflection.Register(s.grpcServer)

	signal.Notify(s.SignalChan(), syscall.SIGTERM, syscall.SIGINT)

	return s
}

func (s *Server) Main() {
	err := s.start()
	if err == nil {
		log.Printf("%v is started", s.ServerCore.Conf.ServiceName)
		defer log.Printf("Server exited")
		waitForSignal := <-s.SignalChan()
		log.Printf("* signal %v", waitForSignal)
		s.Shutdown()
	}
}

func (s *Server) start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", s.ServerCore.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s.grpcListener = lis
	log.Printf("server listening at %v", lis.Addr())

	go s.grpcServer.Serve(lis)

	return nil
}

func (s *Server) Shutdown() {
	log.Printf("Shutting down...")
	s.grpcServer.GracefulStop()
}
