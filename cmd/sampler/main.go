package main

import (
	"flag"
	"context"
	"log"
	"os"
	// "bytes"
	"reflect"

	// "github.com/LukaszSwolkien/IngestTools/metric"
	// "github.com/LukaszSwolkien/IngestTools/shared"
	"github.com/LukaszSwolkien/IngestTools/trace"
	"github.com/apache/thrift/lib/go/thrift"
)

var (
	fileName 	= flag.String("file", "payload.data", "filename for data dump")
	format 		= flag.String("f", "", "The request data format (zipkin, otlp, sapm, jaegerthrift)")
)

func dumpJaegerThrift() {
	sample := trace.GenerateJeagerThriftSample()
	if data, err := thrift.NewTSerializer().Write(context.Background(), &sample); err != nil || len(data) == 0 {
		log.Printf("Error during Thrift serialization of type %v: %v", reflect.TypeOf(sample),err)
	} else {
		os.WriteFile(*fileName, data, 0666)
	}
}

func main() {
	flag.Parse()
	log.Printf("filename `%v`, format `%v`", *fileName, *format)
	switch *format {
	case "jaegerthrift":
		dumpJaegerThrift()
	}

}