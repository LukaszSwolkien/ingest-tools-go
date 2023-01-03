package trace

import (
	"math/rand"
	"time"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
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

func SendZipkinTraceSample(url string, secret string, trace ZipkinTrace) {
	contentType := "application/json"
	json_data, err := json.MarshalIndent(trace, "", "\t")
	if err != nil {
		log.Fatalf("Marshal: %v", err)
	}
	log.Println("Sending sample data:\n" + string(json_data))

	r, err := http.NewRequest("POST", url, bytes.NewBuffer(json_data))
	r.Header.Add("Content-Type", contentType)
	r.Header.Add("X-SF-Token", secret)
	if err != nil {
		log.Fatalln(err)
	}

	client := &http.Client{}

	resp, err := client.Do(r)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	dr, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(string(dr))
}