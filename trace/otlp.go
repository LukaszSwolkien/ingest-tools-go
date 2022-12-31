// sample golang generator for trace ingest data in OTLP format
package trace

import (
	"math/rand"
	"math/big"
	"time"
	// common_v1 "go.opentelemetry.io/proto/otlp/common/v1"
	trace_v1 "go.opentelemetry.io/proto/otlp/trace/v1"

)

func newID() []byte {
	r := rand.Int63()
	b := make([]byte, 8)
	return big.NewInt(r).FillBytes(b)
}


func GetOtlpTrace() *trace_v1.Span {
	now := uint64(time.Now().UnixNano())
	rand.Seed(int64(now))

	traceID := append(newID(), newID()...)
	s := generateSpan(traceID, nil, "root", now)


	return s
}

func generateSpan(traceID []byte, parentSpanID []byte, name string, start uint64) *trace_v1.Span {
	spanID := traceID[:8]
	if parentSpanID != nil {
		spanID = newID()
	}
	return &trace_v1.Span{
		TraceId:                traceID,
		SpanId:                 spanID,
		TraceState:             "",
		ParentSpanId:           parentSpanID,
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
