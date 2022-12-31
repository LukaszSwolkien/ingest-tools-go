package trace

import (
	"fmt"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	trace_v1 "go.opentelemetry.io/proto/otlp/trace/v1"
)

func PostTraceSample(url string, secret string, contentType string, trace ZipkinTrace) {
	json_data, err := json.MarshalIndent(trace, "", "\t")
	if err != nil {
		log.Fatalf("Marshal: %v", err)
	}
	log.Println("* Ingest endpoint: " + url)
	log.Println("* Sending sample data:\n" + string(json_data))

	r, err := http.NewRequest("POST", url+"/v2/trace", bytes.NewBuffer(json_data))
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

func GRPCTraceSample(url string, secret string, port int, data *trace_v1.Span) {
	log.Printf("Seting up a gRPC connection to %s:%s...", url, fmt.Sprint(port))
	err := "not implemented"
	log.Fatalf("did not connect: %v", err)
}
