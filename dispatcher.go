package main

import (
	logevent "github.com/LukaszSwolkien/IngestTools/log"
	"github.com/LukaszSwolkien/IngestTools/metric"
	"github.com/LukaszSwolkien/IngestTools/shared"
	"github.com/LukaszSwolkien/IngestTools/trace"
	"log"
	"time"
)

type dispatcherConfig struct {
	ingest       string
	protocol     string
	transport    string
	token        string
	ingestUrl    string
	grpcInsecure bool
}

func setup(conf dispatcherConfig) *dispatcher {
	return &dispatcher{config: conf}
}

type dispatcher struct {
	config dispatcherConfig
}

func (d *dispatcher) log_sample() {
	switch d.config.protocol {
	case "hec":
		spl_event := logevent.GenerateLogSample()
		content_type := "application/json"
		log.Printf("Splunk Event Log format, Content-Type: %v", content_type)
		shared.SendDataSample(d.config.ingestUrl, d.config.token, content_type, spl_event)
	default:
		log.Fatalf("Unsupported protocol %v", d.config.protocol)
	}
}

func (d *dispatcher) metrics_sample() {
	switch d.config.protocol {
	case "sfx":
		sfx_guage := metric.GenerateSfxGuageDatapointSample()
		content_type := "application/json"
		log.Printf("SignalFx Datapoint format, Content-Type: %v", content_type)
		shared.SendDataSample(d.config.ingestUrl, d.config.token, content_type, sfx_guage)
		for i := 0; i < 3; i++ {
			time.Sleep(time.Second)
			sfx_counter := metric.GenerateSfxCounterDatapointSample()
			shared.SendDataSample(d.config.ingestUrl, d.config.token, content_type, sfx_counter)
		}
	case "otlp":
		otlp_metric := metric.GenerateOtlpMetric()
		metric.SendGrpcOtlpMetricSample(*ingest, *token, *grpcInsecure, otlp_metric)
	default:
		log.Fatalf("Unsupported protocol %v", d.config.protocol)
	}
}

func (d *dispatcher) trace_sample() {
	switch d.config.protocol {
	case "zipkin":
		content_type := "application/json"
		log.Printf("Zipkin JSON format, Content-Type: %v", content_type)
		var zipkin_data = trace.GenerateZipkinSample()
		shared.SendDataSample(d.config.ingestUrl, d.config.token, content_type, zipkin_data)
	case "thrift":
		content_type := "x-thrift"
		log.Fatalf("Jaeger Thrift format not implemented, Content-Type: %v", content_type)
	case "sapm":
		content_type := "x-protobuf"
		log.Fatalf("SAPM format not implemented, Content-Type: %v", content_type)
	case "otlp":
		otlpSpan := trace.GenerateSpan()
		trace.SendGrpcOtlpTraceSample(d.config.ingestUrl, d.config.token, *grpcInsecure, otlpSpan)
	default:
		log.Fatalf("Unsupported protocol %v", d.config.protocol)
	}
}

func (d *dispatcher) dispatch() {
	switch d.config.ingest {
	case "log":
		d.log_sample()
	case "metrics":
		d.metrics_sample()
	case "trace":
		d.trace_sample()
	default:
		log.Fatalln("Unsupported ingest type")
	}
}
