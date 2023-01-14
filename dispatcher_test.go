package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/LukaszSwolkien/IngestTools/ut"
	"github.com/LukaszSwolkien/IngestTools/cmd/mock/trace-server/server"
	mock "github.com/LukaszSwolkien/IngestTools/cmd/mock/server"
	"sync"
	"log"
)


var (
	svr = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	tests = []dispatcherConfig{
	{
		ingest:       "logs",
		ingestUrl:    svr.URL,
		token:        "t0p_s3cr3t",
		format:       "hec",
		transport:    "http",
	},
	{
		ingest:       "metrics",
		ingestUrl:    svr.URL,
		token:        "t0p_s3cr3t",
		format:       "sfx",
		transport:    "http",
	},
	{
		ingest:       "metrics",
		ingestUrl:    svr.URL,
		token:        "t0p_s3cr3t",
		format:       "otlp",
		transport:    "http",
	},
	{
		ingest:       "trace",
		ingestUrl:    svr.URL,
		token:        "t0p_s3cr3t",
		format:       "zipkin",
		transport:    "http",
	},
}
)

func TestHttpCalls(t *testing.T) {
	for _, c := range tests {
		d := setup(c)
		resp := d.dispatch()
		ut.AssertTrue(t, resp == 200)
	}
}

func asyncDo(t *testing.T, s *server.Server, f func()int, wg *sync.WaitGroup) {
	defer s.Shutdown()
	defer wg.Done()
	log.Printf("asyncDo")
	resp := f()
	ut.AssertTrue(t, resp == 200)
	log.Printf("assertTrue resp %v", resp)
}

func TestGrpcOtlpTrace(t *testing.T) {
	s := server.New(mock.Conf{
		ServiceName: "trace-server-mock",
		GrpcPort:    uint16(8201),
	})

	d := setup(dispatcherConfig{
		ingest: "trace",
		format: "otlp",
		ingestUrl: "localhost:8201",
		token: "t0p_s3cr3t",
		grpcInsecure: true,
		transport: "grpc",
	})
	var wg sync.WaitGroup
	wg.Add(1)
	go asyncDo(t, s, d.dispatch, &wg)
	log.Printf("starting server")
	go s.Main()
	log.Printf("Waiting")
	wg.Wait()
}