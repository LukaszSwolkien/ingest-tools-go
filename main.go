package main

import (
	"flag"
	"log"
	"github.com/LukaszSwolkien/IngestTools/shared"
	"github.com/LukaszSwolkien/IngestTools/trace"
)

var (
	ingest = flag.String("ingest", "", "ingest url")
	endpoint = flag.String("endpoint", "", "The ingest type (v1/log, v2/trace, v2/trace/otlp, v2/datapoint, v2/datapoint/otlp)")
	token = flag.String("token", "", "Ingest access token")
	protocol = flag.String("protocol", "zipkin", "The request protocol (zipkin, otlp, sapm, thrift)")
	grpcInsecure = flag.Bool("grpc-insecure", false, "Set grpc-insecure=false to enable TLS")
	// transport = flag.String("transport", "grpc", "Transport (http, grpc)")
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
}

func main() {
	flag.Parse()
	loadConfiguration()
	log.Printf("Ingest endpoint: %v", ingest_url)
	switch *endpoint {
		case "v2/trace":
			if (*protocol == "zipkin"){
				content_type := "application/json"
				log.Printf("Zipkin JSON format, Content-Type: %v", content_type)
				var json_data = trace.GenerateZipkinSample()
				trace.SendZipkinTraceSample(ingest_url, *token, json_data)
			} else if (*protocol == "thrift") {
				content_type := "x-thrift"
				log.Fatalf("Jaeger Thrift format not implemented, Content-Type: %v", content_type)
			} else if (*protocol == "sapm") {
				content_type := "x-protobuf"
				log.Fatalf("SAPM format not implemented, Content-Type: %v", content_type)
			}
		case "v2/trace/otlp":
			log.Printf("Protocol: otlp")
			log.Printf("Transport: grpc")
			otlp_data := trace.GenerateOtlpTrace()
			trace.SendGrpcOtlpTraceSample(ingest_url, *token, *grpcInsecure, otlp_data)
	default:
		log.Fatalln("Unsupported endpoint")
	}
}