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

	commands	map[string]func()
}

func addCommand(commands map[string]func(), cmd string, action func()) {
	commands[cmd] = action
}

func dispatchCommand(commands	map[string]func(), cmd string) {
	if f, ok := commands[cmd]; ok {
		f()
	} else {
		log.Fatalln("Unsupported command")
	}
}

func setup(conf dispatcherConfig) *dispatcher {
	d := &dispatcher{config: conf}
	d.commands = make(map[string]func())
	addCommand(d.commands, "log", d.log_sample)
	addCommand(d.commands, "metrics", d.metrics_sample)
	addCommand(d.commands, "trace", d.trace_sample)
	addCommand(d.commands, "rum", d.rum_sample)
	addCommand(d.commands, "events", d.events_sample)
	addCommand(d.commands, "profiling", d.profiling_sample)
	return d
}

func (d *dispatcher) dispatch() {
	dispatchCommand(d.commands, d.config.ingest)
}

func notImplementedSample(conf dispatcherConfig){
	log.Fatalf("Sample for `%v` ingest and `%v` data format is not implemented.", conf.ingest, conf.format)
}

func unsupportedDataFormat(conf dispatcherConfig){
	log.Fatalf("Unsupported data format `%v` for `%v` ingest", conf.format, conf.ingest)
}

// ---- rum samples ----
func (d *dispatcher) rum_sample() {
	// TODO
	notImplementedSample(d.config)
}
// ---- events samples ----
func (d *dispatcher) events_sample() {
	// TODO
	notImplementedSample(d.config)
}
// ---- profiling samples ----
func (d *dispatcher) profiling_sample() {
	// TODO
	notImplementedSample(d.config)
}
// ---- log samples ----
func (d *dispatcher) log_sample() {
	switch d.config.format {
	case "hec":
		log_hec_sample(d.config)
	default:
		unsupportedDataFormat(d.config)
	}
}
func log_hec_sample(conf dispatcherConfig) {
	spl_event := logevent.GenerateLogSample()
	content_type := "application/json"
	log.Printf("Splunk Event Log format, Content-Type: %v", content_type)
	shared.SendData(conf.ingestUrl, conf.token, content_type, spl_event)
}
// ---- metrics samples ----
func (d *dispatcher) metrics_sample() {
	switch d.config.format {
	case "sfx":
		httpSfxMetrics(d.config)
	case "otlp":
		httpOtlpMetrics(d.config)
	default:
		unsupportedDataFormat(d.config)
	}
}
func httpOtlpMetrics(conf dispatcherConfig) {
	otlp_metric := metric.GenerateOtlpMetric()
	metric.PostOtlpMetricSample(conf.ingestUrl, *token, otlp_metric)
}
func httpSfxMetrics(conf dispatcherConfig) {
	sfx_guage := metric.GenerateSfxGuageDatapointSample()
	content_type := "application/json"
	log.Printf("SignalFx Datapoint format, Content-Type: %v", content_type)
	shared.SendData(conf.ingestUrl, conf.token, content_type, sfx_guage)
	for i := 0; i < 3; i++ {
		time.Sleep(time.Second)
		sfx_counter := metric.GenerateSfxCounterDatapointSample()
		shared.SendData(conf.ingestUrl, conf.token, content_type, sfx_counter)
	}
}
// ---- trace samples ----
func (d *dispatcher) trace_sample() {
	switch d.config.format {
	case "zipkin":
		httpZipkinTrace(d.config)
	case "thrift":
		thriftJaegerTrace(d.config)
	case "sapm":
		httpSapmTrace(d.config)
	case "otlp":
		if d.config.transport == "grpc" {
			grpcOtlpTrace(d.config)
		}else {
			log.Fatalf("'%v' data-format over '%v' not implemented", d.config.format, d.config.transport)
		}
	default:
		unsupportedDataFormat(d.config)
	}
}
func grpcOtlpTrace(conf dispatcherConfig) {
		otlpSpan := trace.GenerateSpan()
		trace.SendGrpcOtlpTraceSample(conf.ingestUrl, conf.token, *grpcInsecure, otlpSpan)
} 
func httpSapmTrace(conf dispatcherConfig) {
	content_type := "x-protobuf"
	log.Fatalf("SAPM format not implemented, Content-Type: %v", content_type)
}
func thriftJaegerTrace(conf dispatcherConfig) {
	// TODO: implement Jaeger
	content_type := "x-thrift"
	log.Fatalf("Jaeger Thrift format not implemented, Content-Type: %v", content_type)
}
func httpZipkinTrace(conf dispatcherConfig) {
	content_type := "application/json"
	log.Printf("Zipkin JSON format, Content-Type: %v", content_type)
	var zipkin_data = trace.GenerateZipkinSample()
	shared.SendData(conf.ingestUrl, conf.token, content_type, zipkin_data)
}
