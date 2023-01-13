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
	format       string
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
	switch d.config.format {
	case "hec":
		spl_event := logevent.GenerateLogSample()
		content_type := "application/json"
		log.Printf("Splunk Event Log format, Content-Type: %v", content_type)
		shared.SendData(d.config.ingestUrl, d.config.token, content_type, spl_event)
	default:
		log.Fatalf("Unsupported format %v", d.config.format)
	}
}

func (d *dispatcher) metrics_sample() {
	switch d.config.format {
	case "sfx":
		sfx_guage := metric.GenerateSfxGuageDatapointSample()
		content_type := "application/json"
		log.Printf("SignalFx Datapoint format, Content-Type: %v", content_type)
		shared.SendData(d.config.ingestUrl, d.config.token, content_type, sfx_guage)
		for i := 0; i < 3; i++ {
			time.Sleep(time.Second)
			sfx_counter := metric.GenerateSfxCounterDatapointSample()
			shared.SendData(d.config.ingestUrl, d.config.token, content_type, sfx_counter)
		}
	case "otlp":
		otlp_metric := metric.GenerateOtlpMetric()
		metric.PostOtlpMetricSample(d.config.ingestUrl, *token, otlp_metric)
		// metric.SendGrpcOtlpMetricSample(*ingest, *token, *grpcInsecure, otlp_metric)
	default:
		log.Fatalf("Unsupported format %v", d.config.format)
	}
}

func (d *dispatcher) trace_sample() {
	switch d.config.format {
	case "zipkin":
		content_type := "application/json"
		log.Printf("Zipkin JSON format, Content-Type: %v", content_type)
		var zipkin_data = trace.GenerateZipkinSample()
		shared.SendData(d.config.ingestUrl, d.config.token, content_type, zipkin_data)
	case "thrift":
		content_type := "x-thrift"
		// TODO: implement Jaeger
		log.Fatalf("Jaeger Thrift format not implemented, Content-Type: %v", content_type)
	case "sapm":
		content_type := "x-protobuf"
		// TODO: implement SAPM
		log.Fatalf("SAPM format not implemented, Content-Type: %v", content_type)
	case "otlp":
		if d.config.transport == "grpc" {
			otlpSpan := trace.GenerateSpan()
			trace.SendGrpcOtlpTraceSample(d.config.ingestUrl, d.config.token, *grpcInsecure, otlpSpan)
		} else {
			log.Fatalf("'%v' data-format over '%v' not implemented", d.config.format, d.config.transport)
		}
	default:
		log.Fatalf("Unsupported format %v", d.config.format)
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
