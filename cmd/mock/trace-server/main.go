package main

import (
	"flag"

	core "github.com/LukaszSwolkien/ingest-tools/cmd/mock/server"
	"github.com/LukaszSwolkien/ingest-tools/cmd/mock/trace-server/server"
)

var (
	port = flag.Uint("port", 8201, "server port")
)

func main() {
	flag.Parse()
	s := server.New(core.Conf{
		ServiceName: "trace-server-mock",
		GrpcPort:    uint16(*port),
	})

	s.Main()
}
