package trace

import (
	"os"
	"fmt"
	"bytes"
	"encoding/json"
	"log"
	"crypto/tls"
	"net/http"
	"net/http/httputil"
	"context"

	trace_ingest "go.opentelemetry.io/proto/otlp/collector/trace/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	sfx_grpc_auth "github.com/signalfx/ingest-protocols/grpc"
)

func PostZipkinTraceSample(url string, secret string, trace ZipkinTrace) {
	contentType := "application/json"
	json_data, err := json.MarshalIndent(trace, "", "\t")
	if err != nil {
		log.Fatalf("Marshal: %v", err)
	}
	log.Println("Sending sample data:\n" + string(json_data))

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	r.Header.Add("Content-Type", contentType)
	r.Header.Add("X-SF-Token", secret)
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}

	resp, err := client.Do(r)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	dr, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(dr))
}

func GrpcOtlpTraceSample(url string, secret string, grpcInsecure bool, data  *trace_ingest.ExportTraceServiceRequest) {
	var sec grpc.DialOption
	security_option := "insecure"
	if grpcInsecure {
		sec = grpc.WithTransportCredentials(insecure.NewCredentials())
	} else {
		sec = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))
		security_option = "TLS"
	}
	log.Printf("Security option: %v", security_option)
	log.Printf("Setting up a gRPC connection to %v", fmt.Sprint(url))
	conn, err := grpc.Dial(url, sec, grpc.WithBlock())
	if err != nil {
		log.Printf("Failed to connect to gRPC server: %v", err)
		os.Exit(1)
	}
	log.Printf("Connection successful %v", conn)

	auth := &sfx_grpc_auth.SignalFxTokenAuth{Token: secret, DisableTransportSecurity: grpcInsecure}

	ingest_cli := trace_ingest.NewTraceServiceClient(conn)
	rs, err := ingest_cli.Export(context.Background(), data, grpc.PerRPCCredentials(auth))
	if err != nil {
		log.Printf("Failed to call gRPC Export methond: %v", err)
		os.Exit(1)
	}
	log.Printf("resoponse %v", rs.String())
}
