package main

import (
	"fmt"
	"log"

	logevent "github.com/LukaszSwolkien/IngestTools/log"
	"github.com/LukaszSwolkien/IngestTools/metric"
	"github.com/LukaszSwolkien/IngestTools/shared"
	"github.com/LukaszSwolkien/IngestTools/trace"
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

	commands map[string]func() int
}

func addCommand(commands map[string]func() int, cmd string, action func() int) {
	commands[cmd] = action
}

func dispatchCommand(commands map[string]func() int, cmd string) int {
	if f, ok := commands[cmd]; ok {
		return f()
	} else {
		log.Printf("Unsupported command")
		return 400
	}
}

func setup(conf dispatcherConfig) *dispatcher {
	d := &dispatcher{config: conf}
	d.commands = make(map[string]func() int)
	addCommand(d.commands, "logs", d.logsSample)
	addCommand(d.commands, "metrics", d.metricsSample)
	addCommand(d.commands, "trace", d.traceSample)
	// TODO: implement samples for the following:
	// addCommand(d.commands, "rum", d.rumSample)
	// addCommand(d.commands, "events", d.eventsSample)
	// addCommand(d.commands, "profiling", d.profilingSample)
	return d
}

func (d *dispatcher) dispatch() int {
	return dispatchCommand(d.commands, d.config.ingest)
}

func unsupportedDataFormat(conf dispatcherConfig) string {
	return fmt.Sprintf("Unsupported data format `%v` for `%v` ingest", conf.format, conf.ingest)
}

// ---- log samples ----
func (d *dispatcher) logsSample() int {
	switch d.config.format {
	case "hec":
		return d.httpHeclogs()
	default:
		log.Printf(unsupportedDataFormat(d.config))
	}
	return 400
}
func (d *dispatcher) httpHeclogs() int {
	spl_event := logevent.GenerateLogSample()
	content_type := "application/json"
	log.Printf("Splunk Event Log format, Content-Type: %v", content_type)
	return shared.SendJsonData(d.config.ingestUrl, d.config.token, content_type, spl_event)
}

// ---- metrics samples ----
func (d *dispatcher) httpMetricsSample() int {
	switch d.config.format {
	case "sfx":
		return d.httpSfxMetrics()
	case "otlp":
		return d.httpOtlpMetrics()
	default:
		log.Printf(unsupportedDataFormat(d.config))
	}
	return 400
}

func (d *dispatcher) metricsSample() int {
	switch d.config.transport {
	case "http":
		return d.httpMetricsSample()
	default:
		log.Printf(unsupportedDataFormat(d.config))
	}
	return 400
}
func (d *dispatcher) httpOtlpMetrics() int {
	otlp_metric := metric.GenerateOtlpMetric()
	return metric.PostOtlpMetricSample(d.config.ingestUrl, *token, otlp_metric)
}
func (d *dispatcher) httpSfxMetrics() int {
	sfx_guage := metric.GenerateSfxGuageDatapointSample()
	content_type := "application/json"
	log.Printf("SignalFx Datapoint format, Content-Type: %v", content_type)
	r := shared.SendJsonData(d.config.ingestUrl, d.config.token, content_type, sfx_guage)
	// for i := 0; i < 3; i++ {
	// 	time.Sleep(time.Second)
	// 	sfx_counter := metric.GenerateSfxCounterDatapointSample()
	// 	r = shared.SendData(d.config.ingestUrl, d.config.token, content_type, sfx_counter)
	// }
	return r
}

func (d *dispatcher) httpTraceSample() int {
	switch d.config.format {
	case "zipkin":
		return d.httpZipkinTrace()
	// case "thrift":
	// 	return d.thriftJaegerTrace()
	// case "sapm":
	// 	return d.httpSapmTrace()
	case "otlp":
		return d.httpOtlpTrace()
	// log.Printf("'%v' data-format over '%v' not implemented", d.config.format, d.config.transport)
	default:
		log.Printf(unsupportedDataFormat(d.config))
	}
	return 400
}

func (d *dispatcher) grpcTraceSample() int {
	switch d.config.format {
	case "otlp":
		return d.grpcOtlpTrace()
	default:
		log.Printf(unsupportedDataFormat(d.config))
	}
	return 400
}

// ---- trace samples ----
func (d *dispatcher) traceSample() int {
	switch d.config.transport {
	case "http":
		return d.httpTraceSample()
	case "grpc":
		return d.grpcTraceSample()
	default:
		log.Printf(unsupportedDataFormat(d.config))
	}
	return 400
}
func (d *dispatcher) grpcOtlpTrace() int {
	otlpSpan := trace.GenerateSpan()
	return trace.SendGrpcOtlpTraceSample(d.config.ingestUrl, d.config.token, d.config.grpcInsecure, otlpSpan)
}

func (d *dispatcher) httpZipkinTrace() int {
	content_type := "application/json"
	log.Printf("Zipkin JSON format, Content-Type: %v", content_type)
	var zipkin_data = trace.GenerateZipkinSample()
	return shared.SendJsonData(d.config.ingestUrl, d.config.token, content_type, zipkin_data)
}

func (d *dispatcher) httpOtlpTrace() int {
	content_type := "application/x-protobuf"
	log.Printf("OTLP protobuf format, Content-Type: %v", content_type)
	otlpSpan := trace.GenerateSpan()
	return trace.SendHttpOtlpSample(d.config.ingestUrl, d.config.token, content_type, otlpSpan)
}
