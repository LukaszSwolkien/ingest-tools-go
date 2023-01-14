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

type dispatcher struct {
	config dispatcherConfig

	commands	map[string]func()int
}

func addCommand(commands map[string]func()int, cmd string, action func()int) {
	commands[cmd] = action
}

func dispatchCommand(commands	map[string]func()int, cmd string) int {
	if f, ok := commands[cmd]; ok {
		return f()
	} else {
		log.Fatalln("Unsupported command")
		return 404
	}
}

func setup(conf dispatcherConfig) *dispatcher {
	d := &dispatcher{config: conf}
	d.commands = make(map[string]func()int)
	addCommand(d.commands, "logs", d.logs_sample)
	addCommand(d.commands, "metrics", d.metrics_sample)
	addCommand(d.commands, "trace", d.trace_sample)
	addCommand(d.commands, "rum", d.rum_sample)
	addCommand(d.commands, "events", d.events_sample)
	addCommand(d.commands, "profiling", d.profiling_sample)
	return d
}

func (d *dispatcher) dispatch() int {
	return dispatchCommand(d.commands, d.config.ingest)
}

func notImplementedSample(conf dispatcherConfig){
	log.Fatalf("Sample for `%v` ingest and `%v` data format is not implemented.", conf.ingest, conf.format)
}

func unsupportedDataFormat(conf dispatcherConfig){
	log.Fatalf("Unsupported data format `%v` for `%v` ingest", conf.format, conf.ingest)
}

// ---- rum samples ----
func (d *dispatcher) rum_sample() int {
	// TODO
	notImplementedSample(d.config)
	return 404
}
// ---- events samples ----
func (d *dispatcher) events_sample() int{
	// TODO
	notImplementedSample(d.config)
	return 404

}
// ---- profiling samples ----
func (d *dispatcher) profiling_sample() int{
	// TODO
	notImplementedSample(d.config)
	return 404
}
// ---- log samples ----
func (d *dispatcher) logs_sample() int{
	switch d.config.format {
	case "hec":
		return httpHeclogs(d.config)
	default:
		unsupportedDataFormat(d.config)
	}
	return 404
}
func httpHeclogs(conf dispatcherConfig) int {
	spl_event := logevent.GenerateLogSample()
	content_type := "application/json"
	log.Printf("Splunk Event Log format, Content-Type: %v", content_type)
	return shared.SendData(conf.ingestUrl, conf.token, content_type, spl_event)
}
// ---- metrics samples ----
func (d *dispatcher) metrics_sample() int {
	switch d.config.format {
	case "sfx":
		return httpSfxMetrics(d.config)
	case "otlp":
		return httpOtlpMetrics(d.config)
	default:
		unsupportedDataFormat(d.config)
	}
	return 404
}
func httpOtlpMetrics(conf dispatcherConfig) int {
	otlp_metric := metric.GenerateOtlpMetric()
	return metric.PostOtlpMetricSample(conf.ingestUrl, *token, otlp_metric)
}
func httpSfxMetrics(conf dispatcherConfig) int {
	sfx_guage := metric.GenerateSfxGuageDatapointSample()
	content_type := "application/json"
	log.Printf("SignalFx Datapoint format, Content-Type: %v", content_type)
	r := shared.SendData(conf.ingestUrl, conf.token, content_type, sfx_guage)
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		sfx_counter := metric.GenerateSfxCounterDatapointSample()
		r = shared.SendData(conf.ingestUrl, conf.token, content_type, sfx_counter)
	}
	return r
}
// ---- trace samples ----
func (d *dispatcher) trace_sample() int {
	switch d.config.format {
	case "zipkin":
		return httpZipkinTrace(d.config)
	// case "thrift":
	// 	return thriftJaegerTrace(d.config)
	// case "sapm":
	// 	return httpSapmTrace(d.config)
	case "otlp":
		if d.config.transport == "grpc" {
			return grpcOtlpTrace(d.config)
		}else {
			log.Fatalf("'%v' data-format over '%v' not implemented", d.config.format, d.config.transport)
		}
	default:
		unsupportedDataFormat(d.config)
	}
	return 404
}
func grpcOtlpTrace(conf dispatcherConfig) int {
		otlpSpan := trace.GenerateSpan()
		return trace.SendGrpcOtlpTraceSample(conf.ingestUrl, conf.token, conf.grpcInsecure, otlpSpan)
} 
func httpSapmTrace(conf dispatcherConfig) int {
	content_type := "x-protobuf"
	log.Fatalf("SAPM format not implemented, Content-Type: %v", content_type)
	return 501
}
func thriftJaegerTrace(conf dispatcherConfig) int {
	// TODO: implement Jaeger
	content_type := "x-thrift"
	log.Fatalf("Jaeger Thrift format not implemented, Content-Type: %v", content_type)
	return 501
}
func httpZipkinTrace(conf dispatcherConfig) int {
	content_type := "application/json"
	log.Printf("Zipkin JSON format, Content-Type: %v", content_type)
	var zipkin_data = trace.GenerateZipkinSample()
	return shared.SendData(conf.ingestUrl, conf.token, content_type, zipkin_data)
}
