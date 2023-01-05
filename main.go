package main

import (
	"flag"
	"log"
	"github.com/LukaszSwolkien/IngestTools/shared"
)

var (
	ingest = flag.String("ingest", "", "ingest url")
	endpoint = flag.String("endpoint", "", "The ingest type (v1/log, v2/trace, v2/trace/otlp, v2/datapoint, v2/datapoint/otlp)")
	token = flag.String("token", "", "Ingest access token")
	protocol = flag.String("protocol", "zipkin", "The request protocol (zipkin, otlp, sapm, thrift)")
	grpcInsecure = flag.Bool("grpc-insecure", false, "Set grpc-insecure=false to enable TLS")
	transport = flag.String("transport", "grpc", "Transport (http, grpc)")
	ingest_url = ""
)

func loadConfiguration(){
	var c shared.Conf
	err := c.LoadConf(".conf.yaml")
	if err == nil {
		if *ingest == "" {
			*ingest = c.Ingest
		}
		if *token == "" {
			*token = c.Token
		}
		if *endpoint == "" {
			*endpoint = c.Endpoint
		}
		if *protocol == "" {
			*protocol = c.Protocol
		}
	}
	ingest_url = *ingest + "/" + *endpoint
	log.Printf("Ingest endpoint: %v", ingest_url)
	log.Printf("Token: %v", *token)
}


func main() {
	flag.Parse()
	loadConfiguration()
	d := setup(dispatcherConfig{
		ingest: *ingest,
		endpoint: *endpoint,
		token: *token,
		protocol: *protocol,
		grpcInsecure: *grpcInsecure,
		transport: *transport,
	})
	d.dispatch()
}