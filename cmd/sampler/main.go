package main

import (
	"flag"
	"context"
	"log"
	"os"
	// "bytes"
	"reflect"

	// "github.com/LukaszSwolkien/IngestTools/metric"
	"github.com/LukaszSwolkien/IngestTools/shared"
	"github.com/LukaszSwolkien/IngestTools/trace"
	"github.com/apache/thrift/lib/go/thrift"
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
		return dumpJaegerThrift()
	case "otlp":

	default:
		log.Printf("Unsupported data format `%v` for `%v` ingest", *format, *ingest)
	}
	return 0
}

func dumpMetricsSample() int {
	switch *format {
	case "otlp":

	default:
		log.Printf("Unsupported data format `%v` for `%v` ingest", *format, *ingest)
	}
	return 0
}

func dumpJaegerThrift() int {
	sample := trace.GenerateJeagerThriftSample()
	if data, err := thrift.NewTSerializer().Write(context.Background(), &sample); err != nil || len(data) == 0 {
		log.Printf("Error during Thrift serialization of type %v: %v", reflect.TypeOf(sample),err)
		return 0
	} else {
		os.WriteFile(*fileName, data, 0666)
	}
	return 1
}

func main() {
	flag.Parse()
	log.Printf("filename `%v`, format `%v`", *fileName, *format)

	shared.DispatchCommand(commands, *ingest)
}