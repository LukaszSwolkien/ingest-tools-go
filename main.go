package main

import (
	"log"
	"github.com/LukaszSwolkien/IngestTools/shared"
	"github.com/LukaszSwolkien/IngestTools/trace"
)




func main() {
	var c shared.Conf
	c.LoadConf(".secrets.yaml")
	var json_data = trace.GenerateZipkinSample()
	trace.PostTraceSample(c.IngestEndpoint, c.IngestToken, "application/json", json_data)
	log.Println("Done.")
}