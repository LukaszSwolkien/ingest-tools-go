package trace

import (
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/jaegertracing/jaeger/model"
	splunksapm "github.com/signalfx/sapm-proto/gen"

	"bytes"
	"log"

	"github.com/LukaszSwolkien/ingest-tools/shared"
)

func GenerateSapmSpan() *model.Batch {
	batch := &model.Batch{
		Process: &model.Process{ServiceName: "sample-jeager-trace-generator"},
		Spans: []*model.Span{
			{
				TraceID:       model.NewTraceID(0, 1),
				SpanID:        model.NewSpanID(2),
				OperationName: "foo", Duration: time.Microsecond * 1,
				Tags: model.KeyValues{
					{
						Key:   "span.kind",
						VStr:  "client",
						VType: model.StringType,
					},
				},
				StartTime: time.Now().UTC(),
			},
		},
	}
	return batch
}

func SendHttpSapmSample(url string, secret string, contentType string, batch *model.Batch) int {
	bb, err := proto.Marshal(&splunksapm.PostSpansRequest{Batches: []*model.Batch{batch}})
	if err != nil {
		log.Printf("Marshal: %v", err)
		return 400
	}
	body := bytes.NewBuffer(bb)
	return shared.PostHttpRequest(url, secret, contentType, body)
}
