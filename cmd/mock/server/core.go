package server

import "os"

type Conf struct {
	ServiceName string
	HttpPort    uint16
	GrpcPort    uint16
}

type ServerCore struct {
	Conf
	signalChan   chan os.Signal
	// TODO: add anything common to futher send data for debugging
}

// Init the core with the provided config.
func (s *ServerCore) Init(conf Conf) {
	s.Conf = conf
	s.signalChan = make(chan os.Signal)
}

// Return os.Signal to the caller i.e. to callback on the SIGINT etc.
func (s *ServerCore) SignalChan() chan os.Signal {
	return s.signalChan
}
