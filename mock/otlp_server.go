package main

import (
	"fmt"
	"log"
	"context"
	"flag"
	"net"
	status "google.golang.org/grpc/status"
	codes "google.golang.org/grpc/codes"
	colTrace "go.opentelemetry.io/proto/otlp/collector/trace/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 8201, "server port")
)

type OtlpTraceServer struct {
	colTrace.UnimplementedTraceServiceServer
}

func (s *OtlpTraceServer) Export(ctx context.Context, in *colTrace.ExportTraceServiceRequest) (*colTrace.ExportTraceServiceResponse, error) {
	log.Printf("Received: %v", in.String())

	return nil, status.Errorf(codes.Unimplemented, "method Export not implemented")
}

func main(){
	flag.Parse()
	log.Printf("Starting gRPC OTLP server mock")
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	colTrace.RegisterTraceServiceServer(s, &OtlpTraceServer{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}