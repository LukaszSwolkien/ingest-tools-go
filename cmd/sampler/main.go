package main

import (
	"context"
	"flag"
	"log"
	"os"

	"reflect"

	"github.com/LukaszSwolkien/IngestTools/metric"
	"github.com/LukaszSwolkien/IngestTools/shared"
	"github.com/LukaszSwolkien/IngestTools/trace"
	"github.com/apache/thrift/lib/go/thrift"
	"github.com/golang/protobuf/proto"

	traceSvc "go.opentelemetry.io/proto/otlp/collector/trace/v1" // OTLP trace service
	metricSvc "go.opentelemetry.io/proto/otlp/collector/metrics/v1" // OTLP metrics service
)

var (
	fileName 	= flag.String("file", "payload.data", "filename for data dump")
	format 		= flag.String("f", "", "The request data format (zipkin, otlp, sapm, jaegerthrift)")
	ingest      = flag.String("i", "", "ingest type (trace, metrics, logs, events, rum...)")
	commands	= make(map[string]func()int)
)

func init() {
	shared.AddCommand(commands, "trace", dumpTraceSample)	
	shared.AddCommand(commands, "metrics", dumpMetricsSample)	
}

func dumpTraceSample() int {
	switch *format {
	case "jaegerthrift":
		return dumpTraceJaegerThrift()
	case "otlp":
		return dumpTraceOtlp()
	default:
		log.Printf("Unsupported data format `%v` for `%v` ingest", *format, *ingest)
	}
	return 0
}

func dumpMetricsSample() int {
	switch *format {
	case "otlp":
		return dumpMetricsOtlp()

	default:
		log.Printf("Unsupported data format `%v` for `%v` ingest", *format, *ingest)
	}
	return 0
}

func dumpMetricsOtlp() int {
	otlp_metric := metric.GenerateOtlpMetric()
	d := &metricSvc.ExportMetricsServiceRequest{ResourceMetrics: metric.GetResourceMetric(otlp_metric)}

	b, err := proto.Marshal(d)
	if err != nil {
		log.Fatalf("Marshal: %v", err)
		return 400
	}
	err = os.WriteFile(*fileName, b, 0666)
	if err != nil {
        log.Printf("Cannot write binary data to file: %v", err)
		return 400
    }

	return 0
}

func dumpTraceOtlp() int {
	sample := trace.GenerateOtlpSpan()
	resSpans := trace.GetResourceSpans(sample)
	data := &traceSvc.ExportTraceServiceRequest{ResourceSpans: resSpans}
	message, err := proto.Marshal(data)
	if err != nil {
		log.Printf("Marshal: %v", err)
		return 400
	}

	err = os.WriteFile(*fileName, message, 0666)
	if err != nil {
        log.Printf("Cannot write binary data to file: %v", err)
		return 400
    }

	return 0
}

func dumpTraceJaegerThrift() int {
	sample := trace.GenerateJeagerThriftSample()
	if data, err := thrift.NewTSerializer().Write(context.Background(), &sample); err != nil || len(data) == 0 {
		log.Printf("Error during Thrift serialization of type %v: %v", reflect.TypeOf(sample),err)
		return 400
	} else {
		os.WriteFile(*fileName, data, 0666)
	}
	return 0
}

func main() {
	flag.Parse()
	log.Printf("filename `%v`, format `%v`", *fileName, *format)

	shared.DispatchCommand(commands, *ingest)
}