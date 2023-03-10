package shared

import (
	yaml "gopkg.in/yaml.v3"
	"log"
	"os"
)

type Conf struct {
	Ingest       string `yaml:"ingest"`
	Format       string `yaml:"data-format"`
	Transport    string `yaml:"transport"`
	Token        string `yaml:"token"`
	IngestUrl    string `yaml:"url"`
	GrpcInsecure string `yaml:"grpc-insecure"`
}

func (c *Conf) LoadConf(confFile string) error {
	yamlFile, err := os.ReadFile(confFile)
	if err != nil {
		log.Printf("cannot read configuration, err to %v", err)
		return err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return err
	}

	return nil
}
