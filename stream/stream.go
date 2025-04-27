package stream

import "log"

type Stream struct {
	clientIP, clientPort, serverIP, serverPort, protocol string
}

func New(clientIP, serverIP, clientPort, serverPort, protocol string) *Stream {
	stream := &Stream{
		clientIP:   clientIP,
		clientPort: clientPort,
		serverIP:   serverIP,
		serverPort: serverPort,
		protocol:   protocol,
	}
	return stream
}

func (s *Stream) Start() {
	switch s.protocol {
	case "udp":
		s.StartUDP()
	case "tcp":
		s.StartTCP()
	default:
		log.Fatal("bad protocol")
	}
}
