package server

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	// status "google.golang.org/grpc/status"
	// codes "google.golang.org/grpc/codes"
	colTrace "go.opentelemetry.io/proto/otlp/collector/trace/v1"

	"github.com/LukaszSwolkien/IngestTools/cmd/mock/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	colTrace.UnimplementedTraceServiceServer
	server.Core
	signalChan   chan os.Signal
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
		signalChan: make(chan os.Signal),
		ctx:        context.Background(),
		grpcServer: grpc.NewServer(),
	}
	s.Core.Init(conf)

	colTrace.RegisterTraceServiceServer(s.grpcServer, &Server{})

	// Register reflection service on gRPC server.
	reflection.Register(s.grpcServer)
	return s
}

// Return os.Signal to the caller i.e. to callback on the SIGINT etc.
func (s *Server) SignalChan() chan os.Signal {
	return s.signalChan
}

func (s *Server) Main() {
	err := s.start()
	if err == nil {
		log.Printf("%v is started", s.Core.Conf.ServiceName)
		defer log.Printf("Server exited")
		waitForSignal := <-s.signalChan
		log.Printf("* signal %v", waitForSignal)
		s.shutdown()
	}
}

func (s *Server) start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", s.Core.GrpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s.grpcListener = lis
	log.Printf("server listening at %v", lis.Addr())

	go s.grpcServer.Serve(lis)

	return nil
}

func (s *Server) shutdown() {
	log.Fatalf("Shuting down...")
	s.grpcServer.GracefulStop()
}
