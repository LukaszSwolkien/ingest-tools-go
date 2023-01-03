// The sample OTLP generator for trace ingest data
package trace

import (
	"log"
	"math/rand"
	"time"
	"context"
	"crypto/tls"
	
	common "go.opentelemetry.io/proto/otlp/common/v1"        		// OTEL commons (basic data representation)
	resource "go.opentelemetry.io/proto/otlp/resource/v1"    		// OTEL resource (metadata)
	trace "go.opentelemetry.io/proto/otlp/trace/v1"                	// OTLP traces data representation
	colTrace "go.opentelemetry.io/proto/otlp/collector/trace/v1"   	// OTLP trace service to push spans
	
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	grpcSfxAuth "github.com/signalfx/ingest-protocols/grpc"
)


func generateSpan() *trace.Span {
	now := uint64(time.Now().UnixNano())

	start := now
	return &trace.Span{
		Name:                   "test-span", // An operation name
		StartTimeUnixNano:      start,       // start timestamp
		EndTimeUnixNano:        start + uint64(rand.Int63n(10000000)+10), // end timestamp
		Attributes:             nil, // A list of key-value pairs
		Events:                 []*trace.Span_Event{}, // A set of zero or more Events
		ParentSpanId:           nil, // Parent's Span identifier
		Links:                  []*trace.Span_Link{}, // Links to zero or more causally-related Spans (via the SpanContext to those related Spans)
		// SpanContext. All the info that identifies Span in the Trace.
		TraceId:                []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
		SpanId:                 []byte{0, 0, 0, 0, 0, 0, 0, 2},
		TraceState:             "",
		Kind:                   trace.Span_SPAN_KIND_CLIENT,
		DroppedAttributesCount: 0,
		DroppedEventsCount:     0,
		DroppedLinksCount:      0,
		Status:                 &trace.Status{},
	}
}

func getResource() *resource.Resource {
	return &resource.Resource{
			Attributes: []*common.KeyValue{
				{
					Key: "service.name",
					Value: &common.AnyValue {
						Value: &common.AnyValue_StringValue{
							StringValue: "otlp-trace-generator",
						},
					},
				},
			},
			DroppedAttributesCount: 0,
		}
}

func getInstrumentationScope() *common.InstrumentationScope{
	return &common.InstrumentationScope{
		Name:    "Ingest-Tools-GO OTLP over gRPC sample trace generator",
		Version: "1.0.0",
	}
}

// A collection of Spans produced by an InstrumentationScope
func getScopeSpans() []*trace.ScopeSpans{
	return []*trace.ScopeSpans {
		{
			Scope: getInstrumentationScope(),		// can be nil
			Spans: []*trace.Span{generateSpan()},
		},
	}
}

func GenerateOtlpTrace() []*trace.ResourceSpans {
	return []*trace.ResourceSpans{
		{
			Resource: getResource(),
			ScopeSpans: getScopeSpans(),
		},
	}
}

func SendGrpcOtlpTraceSample(url string, secret string, grpcInsecure bool, resourceSpans []*trace.ResourceSpans) {
	data := &colTrace.ExportTraceServiceRequest{ResourceSpans: resourceSpans}
	security_option := "insecure"
	var sec grpc.DialOption
	log.Printf("insecure %v", grpcInsecure)
	if grpcInsecure {
		sec = grpc.WithTransportCredentials(insecure.NewCredentials())
	} else {
		sec = grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{}))
		security_option = "TLS"
	}
	log.Printf("Security option: %v", security_option)
	log.Printf("Setting up a gRPC connection to %v", url)
	conn, err := grpc.Dial(url, sec)//, grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	auth := &grpcSfxAuth.SignalFxTokenAuth{Token: secret, DisableTransportSecurity: grpcInsecure}
	c := colTrace.NewTraceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	rs, err := c.Export(ctx, data, grpc.PerRPCCredentials(auth))
	if err != nil {
		log.Fatalf("Err in TraseServiceServer.Export methond: %v", err)
	}
	log.Printf("resoponse %v", rs.String())
}