package metric

import (
	"bytes"
	"context"
	"log"
	"time"

	metricSvc "go.opentelemetry.io/proto/otlp/collector/metrics/v1" // OTLP metrics service
	common "go.opentelemetry.io/proto/otlp/common/v1"
	metric "go.opentelemetry.io/proto/otlp/metrics/v1" // OTLP metrics data representation

	"github.com/LukaszSwolkien/ingest-tools/shared"
	grpcSfxAuth "github.com/signalfx/ingest-protocols/grpc"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func makeDatapoint() *metric.NumberDataPoint {
	return &metric.NumberDataPoint{
		TimeUnixNano: uint64(time.Now().UnixNano()),
		Value:        &metric.NumberDataPoint_AsInt{AsInt: 321},
	}
}

func attributes(m map[string]string) []*common.KeyValue {
	ret := make([]*common.KeyValue, 0, len(m))
	for k, v := range m {
		ret = append(ret, &common.KeyValue{
			Key:   k,
			Value: &common.AnyValue{Value: &common.AnyValue_StringValue{StringValue: v}},
		})
	}
	return ret
}

func datapoint() *metric.NumberDataPoint {
	dp := makeDatapoint()
	dp.Attributes = attributes(map[string]string{
		"k0": "v0",
		"k1": "v1",
	})
	return dp
}

func GenerateOtlpMetric() *metric.Metric {
	return &metric.Metric{
		Name: "sample-gen-heartbeat", // unique name of the metric
		Data: &metric.Metric_Gauge{
			Gauge: &metric.Gauge{
				DataPoints: []*metric.NumberDataPoint{
					datapoint(),
				},
			},
		},
	}
}

// A collection of Spans produced by an InstrumentationScope
func getScopeMetrics(metrics *metric.Metric) []*metric.ScopeMetrics {
	return []*metric.ScopeMetrics{
		{
			// Scope:   shared.GetInstrumentationScope("otlp-metric-generator"), // can be nil
			Metrics: []*metric.Metric{metrics},
		},
	}
}

func GetResourceMetric(metrics *metric.Metric) []*metric.ResourceMetrics {
	return []*metric.ResourceMetrics{
		{
			// Resource:     shared.GetResource("otlp-metric-generator"),
			ScopeMetrics: getScopeMetrics(metrics),
		},
	}
}

func SendGrpcOtlpMetricSample(url string, secret string, grpcInsecure bool, metrics *metric.Metric) {
	data := &metricSvc.ExportMetricsServiceRequest{ResourceMetrics: GetResourceMetric(metrics)}
	conn, err := shared.GrpcConnection(url, grpcInsecure)
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	auth := &grpcSfxAuth.SignalFxTokenAuth{Token: secret, DisableTransportSecurity: grpcInsecure}
	c := metricSvc.NewMetricsServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	rs, err := c.Export(ctx, data, grpc.PerRPCCredentials(auth))
	if err != nil {
		log.Fatalf("Server error: %v", err)
	}
	partialSuccess := rs.GetPartialSuccess()
	if partialSuccess != nil {
		log.Printf("Rejected metrics %v", partialSuccess.ErrorMessage)

	} else {
		log.Printf("Request fully accepted")
	}
}

func PostOtlpMetricSample(url string, secret string, metrics *metric.Metric) int {
	// log.Printf("url: %v, secret: %v", url, secret)
	d := &metricSvc.ExportMetricsServiceRequest{ResourceMetrics: GetResourceMetric(metrics)}

	b, err := proto.Marshal(d)
	if err != nil {
		log.Fatalf("Marshal: %v", err)
	}

	data := bytes.NewBuffer(b)
	return shared.PostHttpRequest(url, secret, "application/x-protobuf", data)
}
