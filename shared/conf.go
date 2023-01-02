package shared

import (
	"log"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type Conf struct {
	Token    		string `yaml:"token"`
	Ingest		   	string `yaml:"ingest"`
	Endpoint 		string `yaml:"endpoint"`
	Protocol		string `yaml:"protocol"`
	GrpcInsecure	string `yaml:"grpc-insecure"`
 
}

func (c *Conf) LoadConf(confFile string) (error) {
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
