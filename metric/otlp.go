package metric

import (
	"log"
	"context"
	"time"
	metricSvc "go.opentelemetry.io/proto/otlp/collector/metrics/v1"	// OTLP metrics service
	metric "go.opentelemetry.io/proto/otlp/metrics/v1"        		// OTLP metrics data representation

	"github.com/LukaszSwolkien/IngestTools/shared"
	grpcSfxAuth "github.com/signalfx/ingest-protocols/grpc"
	"google.golang.org/grpc"

)

func GenerateOtlpMetric() *metric.Metric {
	return &metric.Metric{
		Name: "sample-gen-heartbeat", // unique name of the metric
	}
}

// A collection of Spans produced by an InstrumentationScope
func getScopeMetrics(metrics *metric.Metric) []*metric.ScopeMetrics{
	return []*metric.ScopeMetrics {
		{
			Scope: shared.GetInstrumentationScope("otlp-metric-generator"),		// can be nil
			Metrics: []*metric.Metric{metrics},
		},
	}
}

func getResourceMetric(metrics *metric.Metric) []*metric.ResourceMetrics{
	return []*metric.ResourceMetrics{
		{
			Resource: shared.GetResource("otlp-metric-generator"),
			ScopeMetrics: getScopeMetrics(metrics),
		},
	}
}

func SendGrpcOtlpMetricSample(url string, secret string, grpcInsecure bool, metrics *metric.Metric) {
	data := &metricSvc.ExportMetricsServiceRequest{ResourceMetrics: getResourceMetric(metrics)}
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