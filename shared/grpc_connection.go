package shared

import (
	"crypto/tls"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

// Establish grpc connection
func GrpcConnection(url string, grpcInsecure bool) (*grpc.ClientConn, error) {
	security_option := "insecure"
	var sec grpc.DialOption
	log.Printf("insecure %v", grpcInsecure)
	if grpcInsecure {
		sec = grpc.WithTransportCredentials(insecure.NewCredentials())
	} else {
		sec = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))
		security_option = "TLS"
	}
	log.Printf("Security option: %v", security_option)
	log.Printf("Setting up a gRPC connection to %v", url)
	return grpc.Dial(url, sec, grpc.WithBlock())
}
