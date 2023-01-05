// The sample OTLP generator and sender functions for Splunk Observability trace-ingest endpoint
package trace

import (
	"log"
	"math/rand"
	"time"
	"context"
	"google.golang.org/grpc"

	grpcSfxAuth "github.com/signalfx/ingest-protocols/grpc"
	
	"github.com/LukaszSwolkien/IngestTools/shared"

	trace "go.opentelemetry.io/proto/otlp/trace/v1"             	// OTLP traces data representation
	traceSvc "go.opentelemetry.io/proto/otlp/collector/trace/v1"	// OTLP trace service
)


func GenerateSpan() *trace.Span {
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

// A collection of Spans produced by an InstrumentationScope
func GetScopeSpans(span *trace.Span) []*trace.ScopeSpans{
	return []*trace.ScopeSpans {
		{
			Scope: shared.GetInstrumentationScope("otlp-trace-generator"),		// can be nil
			Spans: []*trace.Span{span},
		},
	}
}

func GetResourceSpans(span *trace.Span) []*trace.ResourceSpans {
	return []*trace.ResourceSpans{
		{
			Resource: shared.GetResource("Ingest-Tools-GO OTLP over gRPC sample trace generator"),
			ScopeSpans: GetScopeSpans(span),
		},
	}
}

func SendGrpcOtlpTraceSample(url string, secret string, grpcInsecure bool, span *trace.Span) {
	resSpans := GetResourceSpans(span)
	data := &traceSvc.ExportTraceServiceRequest{ResourceSpans: resSpans}
	conn, err := shared.GrpcConnection(url, grpcInsecure)
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	auth := &grpcSfxAuth.SignalFxTokenAuth{Token: secret, DisableTransportSecurity: grpcInsecure}
	c := traceSvc.NewTraceServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	rs, err := c.Export(ctx, data, grpc.PerRPCCredentials(auth))
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
	partialSuccess := rs.GetPartialSuccess()
	if partialSuccess != nil {
		log.Printf("Rejected spans %v", partialSuccess.ErrorMessage)

	} else {
		log.Printf("Request fully accepted")
	}
}