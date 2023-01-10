package main

import (
	"flag"
	"github.com/LukaszSwolkien/IngestTools/shared"
	"log"
)

var (
	ingest       = flag.String("i", "", "ingest type (trace, metrics, logs, events, rum...)")
	schema       = flag.String("s", "", "The request schema (zipkin, otlp, sapm, thrift)")
	transport    = flag.String("t", "", "Transport (http, grpc)")
	token        = flag.String("token", "", "Ingest access token")
	ingestUrl    = flag.String("url", "", "The URL to ingest endpoint")
	grpcInsecure = flag.Bool("grpc-insecure", false, "Set grpc-insecure=false to enable TLS")
)

func loadConfiguration() {
	var c shared.Conf
	err := c.LoadConf(".conf.yaml")
	if err == nil {
		if *ingest == "" {
			*ingest = c.Ingest
		}
		if *schema == "" {
			*schema = c.Schema
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
	log.Printf("Schema: %v", *schema)
	log.Printf("Token: %v", *token)
}

func main() {
	flag.Parse()
	loadConfiguration()
	d := setup(dispatcherConfig{
		ingest:       *ingest,
		ingestUrl:    *ingestUrl,
		token:        *token,
		schema:       *schema,
		grpcInsecure: *grpcInsecure,
		transport:    *transport,
	})
	d.dispatch()
}
