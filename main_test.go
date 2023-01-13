package main

import (
	"github.com/LukaszSwolkien/IngestTools/shared"
	yaml "gopkg.in/yaml.v3"
	"os"
	"testing"
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

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Error(err)

	}
}

func checkStr(t *testing.T, want string, result string) {
	if want != result {
		t.Errorf("Result is incorrect, got: %v, want: %v.", result, want)

	}
}

func checkBool(t *testing.T, want bool, result bool) {
	if want != result {
		t.Errorf("Result is incorrect, got: %v, want: %v.", result, want)
	}
}

func generateConfigFile(t *testing.T) {
	c := conf

	y, err := yaml.Marshal(c)
	checkErr(t, err)

	err = os.WriteFile(confFile, []byte(y), 0644)
	checkErr(t, err)

}

func cleanup(t *testing.T) {
	err := os.Remove(confFile)
	checkErr(t, err)
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
	checkStr(t, conf.Ingest, *ingest)
	checkStr(t, conf.Format, *format)
	checkStr(t, conf.IngestUrl, *ingestUrl)
	checkStr(t, conf.Token, *token)
	checkStr(t, conf.Transport, *transport)
	checkBool(t, false, *grpcInsecure)

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
	checkStr(t, "overwrite", *ingest)
	checkStr(t, "overwrite", *format)
	checkStr(t, "overwrite", *ingestUrl)
	checkStr(t, "overwrite", *token)
	checkStr(t, "overwrite", *transport)
}
