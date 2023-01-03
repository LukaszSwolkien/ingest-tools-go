package trace

import (
	"math/rand"
	"time"
)

type Endpoint struct {
	ServiceName string `json:"serviceName,omitempty"`
	Ipv4        string `json:"ipv4"`
	Port        int    `json:"port"`
}

type ZipkinSpan struct {
	Name           string            `json:"name"`
	ID             string            `json:"id"`
	TraceID        string            `json:"traceId"`
	ParentID       string            `json:"parentId"`
	Timestamp      int64             `json:"timestamp"`
	Duration       int               `json:"duration"`
	Kind           string            `json:"kind"`
	LocalEndpoint  Endpoint          `json:"localEndpoint"`
	RemoteEndpoint Endpoint          `json:"remoteEndpoint"`
	Tags           map[string]string `json:"tags"`
}

type ZipkinTrace []ZipkinSpan

func GenerateZipkinSample() []ZipkinSpan {
	var trace []ZipkinSpan
	trace = append(trace, ZipkinSpan{
		Name:      "sample-trace-generator",
		ID:        "1",
		TraceID:   "10",
		ParentID:  "100",
		Timestamp: time.Now().Unix(),
		Duration:  rand.Intn(8000) + 2000,
		Kind:      "SERVER",
		LocalEndpoint: Endpoint{
			ServiceName: "backend",
			Ipv4:        "192.168.99.1",
			Port:        3306,
		},
		RemoteEndpoint: Endpoint{
			Ipv4: "172.19.0.2",
			Port: 58648,
		},
		Tags: map[string]string{
			"http.method": "SET",
			"http.path":   "/api",
		},
	})
	return trace
}
