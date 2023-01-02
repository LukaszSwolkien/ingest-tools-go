package main

import (
	"flag"
	"log"
	"github.com/LukaszSwolkien/IngestTools/shared"
	"github.com/LukaszSwolkien/IngestTools/trace"
)

var (
	endpoint = flag.String("endpoint", "v2/trace", 
	"The ingest type (v1/log, v2/trace, v2/trace/otlp, v2/datapoint, v2/datapoint/otlp)")
	content_type = flag.String("content-type", "application/json", "Content-Type (application/json, x-thrift, x-protobuf)")
	protocol = flag.String("protocol", "zipkin", "The request protocol (zipkin, otlp, sapm, thrift)")
	transport = flag.String("transport", "grpc", "Transport (http, grpc)")
	dial_option = flag.Bool("dial-option", false, "Set dial-option=true to enable TLS")
	token = flag.String("ingest-token", "", "Ingest access token")
)


func main() {
	flag.Parse()
	// log.Println(*it, *ct, *body)
	var c shared.Conf
	c.LoadConf(".secrets.yaml")
	ingest_url := c.IngestEndpoint + "/" + *endpoint
	access_token := *token
	if access_token == "" {
		access_token = c.IngestToken
	}
	log.Printf("Ingest endpoint: %v", ingest_url)
	switch *endpoint {
		case "v2/trace":
			if (*content_type == "application/json" || *protocol == "zipkin"){
				log.Printf("Zipkin JSON format")
				var json_data = trace.GenerateZipkinSample()
				trace.PostZipkinTraceSample(ingest_url, access_token, json_data)
			} else if (*content_type == "x-thrift" || *protocol == "thrift") {
				log.Fatalln("Jaeger Thrift format not implemented")
			} else if (*content_type == "x-protobuf" || *protocol == "sapm") {
				log.Fatalln("SAPM format not implemented")
			}
		case "v2/trace/otlp":
			log.Printf("Protocol: otlp")
			log.Printf("Transport: %v", *transport)
			var otlp_data = trace.GetOtlpTrace()
			if *transport == "grpc" {
				trace.GrpcOtlpTraceSample(ingest_url, access_token, *dial_option, otlp_data)
	}
	default:
		log.Fatalln("Unsupported endpoint")
	}
}