package main

import (
	"log"
	"time"
	"github.com/LukaszSwolkien/IngestTools/shared"
	"github.com/LukaszSwolkien/IngestTools/trace"
	"github.com/LukaszSwolkien/IngestTools/metric"
	logevent "github.com/LukaszSwolkien/IngestTools/log"
)

type dispatcherConfig struct {
	ingest	string
	endpoint string
	token string
	protocol string
	grpcInsecure bool
	transport string
}

func setup(conf dispatcherConfig) *dispatcher{
	return &dispatcher{config: conf}
}

type dispatcher struct {
	config dispatcherConfig

}

func (d *dispatcher)dispatch() {

	switch d.config.endpoint {
		case "v1/log":
			spl_event := logevent.GenerateLogSample()
			content_type := "application/json"
			log.Printf("Splunk Event Log format, Content-Type: %v", content_type)
			shared.SendDataSample(ingest_url, *token, content_type, spl_event)
		case "v2/datapoint":
			sfx_guage := metric.GenerateSfxGuageDatapointSample()
			content_type := "application/json"
			log.Printf("SignalFx Datapoint format, Content-Type: %v", content_type)
			shared.SendDataSample(ingest_url, *token, content_type, sfx_guage)
			for i := 0; i < 3; i++{
				time.Sleep(time.Second)
				sfx_counter := metric.GenerateSfxCounterDatapointSample()
				shared.SendDataSample(ingest_url, *token, content_type, sfx_counter)
			}
		case "v2/datapoint/otlp":
			otlp_metric := metric.GenerateOtlpMetric()
			metric.SendGrpcOtlpMetricSample(*ingest, *token, *grpcInsecure, otlp_metric)
		case "v2/trace":
			if (*protocol == "zipkin"){
				content_type := "application/json"
				log.Printf("Zipkin JSON format, Content-Type: %v", content_type)
				var zipkin_data = trace.GenerateZipkinSample()
				shared.SendDataSample(ingest_url, *token, content_type, zipkin_data)
			} else if (*protocol == "thrift") {
				content_type := "x-thrift"
				log.Fatalf("Jaeger Thrift format not implemented, Content-Type: %v", content_type)
			} else if (*protocol == "sapm") {
				content_type := "x-protobuf"
				log.Fatalf("SAPM format not implemented, Content-Type: %v", content_type)
			}
		case "v2/trace/otlp":
			log.Printf("Protocol: otlp, Transport: grpc")
			otlpSpan := trace.GenerateSpan()
			trace.SendGrpcOtlpTraceSample(*ingest, *token, *grpcInsecure, otlpSpan)
	default:
		log.Fatalln("Unsupported endpoint")
	}
}
