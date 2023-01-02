// sample golang generator for trace ingest data in OTLP format
package trace

import (
	"math/rand"
	"math/big"
	"time"
	
	common_v1 "go.opentelemetry.io/proto/otlp/common/v1"
	trace_v1 "go.opentelemetry.io/proto/otlp/trace/v1"
	resource_v1 "go.opentelemetry.io/proto/otlp/resource/v1"
	trace_ingest "go.opentelemetry.io/proto/otlp/collector/trace/v1"

)

func newID() []byte {
	r := rand.Int63()
	b := make([]byte, 8)
	return big.NewInt(r).FillBytes(b)
}


func generateSpan() *trace_v1.Span {
	now := uint64(time.Now().UnixNano())
	rand.Seed(int64(now))

	trace_id := append(newID(), newID()...)
	span_id := newID()
	name := "root"
	start := now
	return &trace_v1.Span{
		TraceId:                trace_id,
		SpanId:                 span_id,
		TraceState:             "",
		ParentSpanId:           nil,
		Name:                   name,
		Kind:                   trace_v1.Span_SPAN_KIND_CLIENT,
		StartTimeUnixNano:      start,
		EndTimeUnixNano:        start + uint64(rand.Int63n(10000000)+10),
		Attributes:             nil,
		DroppedAttributesCount: 0,
		Events:                 []*trace_v1.Span_Event{},
		DroppedEventsCount:     0,
		Links:                  []*trace_v1.Span_Link{},
		DroppedLinksCount:      0,
		Status:                 &trace_v1.Status{},
	}
}

// func GetOtlpTrace() *trace_ingest.ExportTraceServiceRequest {
// 	var spans []*trace_v1.Span
// 	spans = append(spans, generateSpan())
// 	resource_spans := []*trace_v1.ResourceSpans{
// 		{
// 			Resource: &resource_v1.Resource{
// 				Attributes: []*common_v1.KeyValue{
// 					{
// 						Key: "service.name",
// 						Value: &common_v1.AnyValue {
// 							Value: &common_v1.AnyValue_StringValue{
// 								StringValue: "gdi-ingest",
// 							},
// 						},
// 					},
// 				},
// 				DroppedAttributesCount: 0,
// 			},
// 			ScopeSpans: []*trace_v1.ScopeSpans {
// 				{
// 					Scope: &common_v1.InstrumentationScope{
// 						Name:    "gdi-ingest-grpc",
// 						Version: "0.1.0",
// 					},
// 					Spans: spans,
// 				},
// 			},
// 		},
// 	}


// 	return &trace_ingest.ExportTraceServiceRequest{ResourceSpans: resource_spans}
// }

func GetOtlpTrace() *trace_ingest.ExportTraceServiceRequest {
	var spans []*trace_v1.Span
	spans = append(spans, generateSpan())
	resource_spans := []*trace_v1.ResourceSpans{
		{
			Resource: &resource_v1.Resource{
			},
			ScopeSpans: []*trace_v1.ScopeSpans {
				{
					Scope: &common_v1.InstrumentationScope{
						Name:    "otlp-grpc-sample-ingest",
						Version: "1.0.0",
					},
					Spans: spans,
				},
			},
		},
	}


	return &trace_ingest.ExportTraceServiceRequest{ResourceSpans: resource_spans}
}