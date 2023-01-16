package main

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	mock "github.com/LukaszSwolkien/IngestTools/cmd/mock/server"
	"github.com/LukaszSwolkien/IngestTools/cmd/mock/trace-server/server"
	"github.com/LukaszSwolkien/IngestTools/ut"
)

var (
	svr = httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		}))
	testSupportedConf = []dispatcherConfig{
		{
			ingest:    "logs",
			ingestUrl: svr.URL,
			format:    "hec",
			transport: "http",
		},
		{
			ingest:    "metrics",
			ingestUrl: svr.URL,
			format:    "sfx",
			transport: "http",
		},
		{
			ingest:    "metrics",
			ingestUrl: svr.URL,
			format:    "otlp",
			transport: "http",
		},
		{
			ingest:    "trace",
			ingestUrl: svr.URL,
			format:    "zipkin",
			transport: "http",
		},
		{
			ingest:    "trace",
			ingestUrl: svr.URL,
			format:    "otlp",
			transport: "http",
		},
		{
			ingest:    "trace",
			ingestUrl: svr.URL,
			format:    "sfx",
			transport: "http",
		},
		{
			ingest:    "trace",
			ingestUrl: svr.URL,
			format:    "sapm",
			transport: "http",
		},
		{
			ingest:    "trace",
			ingestUrl: svr.URL,
			format:    "jaegerthrift",
			transport: "http",
		},
	}

	testUnsupportedConf = []dispatcherConfig{
		{
			ingest: "unknown",
		},
		{
			ingest:    "metrics",
			transport: "unknown",
		},
		{
			ingest:    "metrics",
			transport: "http",
			format:    "unknown",
		},
		{
			ingest:    "metrics",
			transport: "unknown",
			format:    "sfx",
		},
		{
			ingest:    "logs",
			transport: "unknown",
		},
		{
			ingest:    "logs",
			transport: "http",
			format:    "unknown",
		},
		{
			ingest:    "trace",
			transport: "unknown",
		},
		{
			ingest:    "trace",
			transport: "http",
			format:    "unknown",
		},
		{
			ingest:    "trace",
			transport: "grpc",
			format:    "unknown",
		},
	}
)

func TestUnsupportedSamples(t *testing.T) {
	for _, c := range testUnsupportedConf {
		d := setup(c)
		ret := d.dispatch()
		ut.AssertTrue(t, ret == 400)
	}
}
func TestSupportedHttpSamples(t *testing.T) {
	for _, c := range testSupportedConf {
		d := setup(c)
		resp := d.dispatch()
		log.Printf("%#v", c)
		ut.AssertEqual(t, resp, 200)
	}
}

func TestGrpcOtlpTrace(t *testing.T) {
	s := server.New(mock.Conf{
		ServiceName: "trace-server-mock",
		GrpcPort:    uint16(8201),
	})

	d := setup(dispatcherConfig{
		ingest:       "trace",
		format:       "otlp",
		ingestUrl:    "localhost:8201",
		token:        "t0p_s3cr3t",
		grpcInsecure: true,
		transport:    "grpc",
	})
	go s.Main()
	asyncDo(t, s, d.dispatch)
}

func asyncDo(t *testing.T, s *server.Server, f func() int) {
	defer s.Shutdown()
	resp := f()
	ut.AssertTrue(t, resp == 200)
}
