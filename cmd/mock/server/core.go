package server

type Conf struct {
	ServiceName		string
	HttpPort		uint16
	GrpcPort		uint16
}

type Core struct{
	Conf
	// TODO: add anything common to futher send data for debugging	
}

// Init the core with the provided config.
func (s *Core) Init(conf Conf) {
	s.Conf = conf
}