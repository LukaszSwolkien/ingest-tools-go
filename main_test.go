package main

import (
	"os"
	"testing"

	"github.com/LukaszSwolkien/IngestTools/shared"
	"github.com/LukaszSwolkien/IngestTools/ut"
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

func generateConfigFile(t *testing.T) {
	c := conf

	y, err := yaml.Marshal(c)
	ut.Check(t, err)

	err = os.WriteFile(confFile, []byte(y), 0644)
	ut.Check(t, err)

}

func cleanup(t *testing.T) {
	err := os.Remove(confFile)
	ut.Check(t, err)
	*ingest = ""
	*ingestUrl = ""
	*transport = ""
	*token = ""
	*grpcInsecure = false
}

func TestLoadFullConfiguration(t *testing.T) {
	defer cleanup(t)
	generateConfigFile(t)
	loadConfiguration(confFile)
	ut.AssertEqual(t, conf.Ingest, *ingest)
	ut.AssertEqual(t, conf.Format, *format)
	ut.AssertEqual(t, conf.IngestUrl, *ingestUrl)
	ut.AssertEqual(t, conf.Token, *token)
	ut.AssertEqual(t, conf.Transport, *transport)
	ut.AssertEqual(t, false, *grpcInsecure)
}

func TestLoadConfigurationOverwrite(t *testing.T) {
	defer cleanup(t)
	generateConfigFile(t)
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
