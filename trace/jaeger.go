package trace

import (
	"bytes"
	"context"
	"encoding/json"
	"log"

	"github.com/LukaszSwolkien/ingest-tools/shared"
	"github.com/apache/thrift/lib/go/thrift"
	jaegerThrift "github.com/jaegertracing/jaeger/thrift-gen/jaeger"
)

const (
	JaegerJSON = `
	{
	  "process": {
		"serviceName": "api",
		"tags": [
		  {
			"key": "hostname",
			"vType": "STRING",
			"vStr": "api246-sjc1"
		  },
		  {
			"key": "ip",
			"vType": "STRING",
			"vStr": "10.53.69.61"
		  },
		  {
			"key": "jaeger.version",
			"vType": "STRING",
			"vStr": "golang"
		  }
		]
	  },
	  "spans": [
		{
		  "traceIdLow": 5951113872249657919,
		  "spanId": 6585752,
		  "parentSpanId": 6866147,
		  "operationName": "get",
		  "startTime": 1485467191639875,
		  "duration": 22938,
		  "tags": [
			{
			  "key": "http.url",
			  "vType": "STRING",
			  "vStr": "http://127.0.0.1:15598/client_transactions"
			},
			{
			  "key": "span.kind",
			  "vType": "STRING",
			  "vStr": "server"
			},
			{
			  "key": "peer.port",
			  "vType": "LONG",
			  "vLong": 53931
			},
			{
			  "key": "someBool",
			  "vType": "BOOL",
			  "vBool": true
			},
			{
			  "key": "someDouble",
			  "vType": "DOUBLE",
			  "vDouble": 129.8
			},
			{
			  "key": "peer.service",
			  "vType": "STRING",
			  "vStr": "rtapi"
			},
			{
			  "key": "peer.ipv4",
			  "vType": "LONG",
			  "vLong": 3224716605
			}
		  ],
		  "logs": [
			{
			  "timestamp": 1485467191639875,
			  "fields": [
				{
				  "key": "key1",
				  "vType": "STRING",
				  "vStr": "value1"
				},
				{
				  "key": "key2",
				  "vType": "STRING",
				  "vStr": "value2"
				}
			  ]
			},
			{
			  "timestamp": 1485467191639875,
			  "fields": [
				{
				  "key": "event",
				  "vType": "STRING",
				  "vStr": "nothing"
				}
			  ]
			}
		  ]
		}
	  ]
	}
	`
)

func GenerateJeagerThriftSample() jaegerThrift.Batch {
	// test data copied from https://github.com/jaegertracing/jaeger/blob/master/model/converter/thrift/jaeger/fixtures/thrift_batch_01.json

	var batch jaegerThrift.Batch
	json.Unmarshal([]byte(JaegerJSON), &batch)
	log.Printf("batch data %v type of %T", batch, batch)
	return batch
}

func SendHttpJaegerThriftSample(url string, secret string, contentType string, sample *jaegerThrift.Batch) int {
	bb, err := thrift.NewTSerializer().Write(context.Background(), sample)
	if err != nil {
		log.Printf("Marshal: %v", err)
		return 400
	}
	body := bytes.NewBuffer(bb)
	return shared.PostHttpRequest(url, secret, contentType, body)
}
