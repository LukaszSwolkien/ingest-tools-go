package shared

import (
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Conf struct {
	IngestToken    string `yaml:"splunk-ingest-token"`
	IngestEndpoint string `yaml:"splunk-ingest"`
}

func (c *Conf) LoadConf(confFile string) *Conf {
	yamlFile, err := os.ReadFile(confFile)
	if err != nil {
		log.Printf("cannot read secrets, err to %v", err)
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
