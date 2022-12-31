package main

import (
	"flag"
	"log"
	"github.com/LukaszSwolkien/IngestTools/shared"
	"github.com/LukaszSwolkien/IngestTools/trace"
)

var (
	it = flag.String("it", "trace", "The ingest type (trace, metric, log)")
	ct = flag.String("ct", "application/json", "Content-Type (application/json, x-thrift, x-protobuf)")
	body = flag.String("body", "zipkin", "The request body type")
	port = flag.Int("port", 8201, "The gRPC target server port")
)


func main() {
	flag.Parse()
	// log.Println(*it, *ct, *body)
	var c shared.Conf
	c.LoadConf(".secrets.yaml")
	switch *it {
	case "trace":
		switch *body {
			case "zipkin":
				var json_data = trace.GenerateZipkinSample()
				trace.PostTraceSample(c.IngestEndpoint, c.IngestToken, *ct, json_data)
			case "otlp":
				var otlp_data = trace.GetOtlpTrace()
				trace.GRPCTraceSample(c.IngestEndpoint, c.IngestToken, *port, otlp_data)
		}
	
	default:
		log.Fatalln("Not implemented")
	}
}