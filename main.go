package main

import (
	"flag"
	"log"

	"github.com/LukaszSwolkien/ingest-tools/shared"
)

var (
	ingest       = flag.String("i", "", "ingest type (trace, metrics, logs, events, rum...)")
	format       = flag.String("f", "", "The request data format (zipkin, otlp, sapm, jaegerthrift)")
	transport    = flag.String("t", "", "Transport (http, grpc)")
	token        = flag.String("token", "", "Ingest access token")
	ingestUrl    = flag.String("url", "", "The URL to ingest endpoint")
	grpcInsecure = flag.Bool("grpc-insecure", false, "Set grpc-insecure=false to enable TLS")
)

func loadConfiguration(fileName string) {
	var c shared.Conf
	err := c.LoadConf(fileName)
	if err == nil {
		if *ingest == "" {
			*ingest = c.Ingest
		}
		if *format == "" {
			*format = c.Format
		}
		if *transport == "" {
			*transport = c.Transport
		}
		if *ingestUrl == "" {
			*ingestUrl = c.IngestUrl
		}
		if *token == "" {
			*token = c.Token
		}
	}
	log.Printf("Ingest: %v", *ingest)
	log.Printf("Ingest endpoint: %v", *ingestUrl)
	log.Printf("Format: %v", *format)
	log.Printf("Token: %v", *token)
}

func main() {
	flag.Parse()
	loadConfiguration(".conf.yaml")
	d := setup(dispatcherConfig{
		ingest:       *ingest,
		ingestUrl:    *ingestUrl,
		token:        *token,
		format:       *format,
		grpcInsecure: *grpcInsecure,
		transport:    *transport,
	})
	d.dispatch()
}
