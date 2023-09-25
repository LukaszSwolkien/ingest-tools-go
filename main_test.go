package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/LukaszSwolkien/ingest-tools/shared"
	"github.com/LukaszSwolkien/ingest-tools/ut"
	yaml "gopkg.in/yaml.v3"
)

const (
	confFile = "test.conf.yaml"
)

var (
	conf = shared.Conf{
		Ingest:       "metrics",
		Format:       "otlp",
		Transport:    "http",
		Token:        "T3StingAcCESt0K3nO11y",
		IngestUrl:    "https://ingest.lab0.signalfx.com",
		GrpcInsecure: "false",
	}
)

func generateConfigFile() {
	c := conf
	y, _ := yaml.Marshal(c)
	err := os.WriteFile(confFile, []byte(y), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func cleanup() {
	*ingest = ""
	*ingestUrl = ""
	*transport = ""
	*token = ""
	*grpcInsecure = false
}

func TestMain(m *testing.M) {
	generateConfigFile()
	exitVal := m.Run()
	err := os.Remove(confFile)
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(exitVal)
}

func TestLoadFullConfiguration(t *testing.T) {
	defer cleanup()
	loadConfiguration(confFile)
	ut.AssertEqual(t, conf.Ingest, *ingest)
	ut.AssertEqual(t, conf.Format, *format)
	ut.AssertEqual(t, conf.IngestUrl, *ingestUrl)
	ut.AssertEqual(t, conf.Token, *token)
	ut.AssertEqual(t, conf.Transport, *transport)
	ut.AssertEqual(t, false, *grpcInsecure)
}

func TestLoadConfigurationOverwrite(t *testing.T) {
	defer cleanup()
	*ingest = "overwrite"
	*format = "overwrite"
	*ingestUrl = "overwrite"
	*token = "overwrite"
	*transport = "overwrite"
	loadConfiguration(confFile)
	ut.AssertEqual(t, "overwrite", *ingest)
	ut.AssertEqual(t, "overwrite", *format)
	ut.AssertEqual(t, "overwrite", *ingestUrl)
	ut.AssertEqual(t, "overwrite", *token)
	ut.AssertEqual(t, "overwrite", *transport)
}

func TestMainFunc(t *testing.T) {
	os.Args = append(os.Args, "-i=metrics")
	os.Args = append(os.Args, "-f=sfx")
	os.Args = append(os.Args, "-t=http")
	os.Args = append(os.Args, fmt.Sprintf("-url=%v", svr.URL))
	os.Args = append(os.Args, "-token=TOKEN")
	main()
}
