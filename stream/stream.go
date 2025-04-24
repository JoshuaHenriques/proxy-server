package stream

import "log"

type Stream struct {
	ClientIP, ClientPort, ServerIP, ServerPort, Protocol string
}

func New(clientIP, serverIP, clientPort, serverPort, protocol string) *Stream {
	stream := &Stream{
		ClientIP:   clientIP,
		ClientPort: clientPort,
		ServerIP:   serverIP,
		ServerPort: serverPort,
		Protocol:   protocol,
	}
	return stream
}

func (s *Stream) Start() {
	switch s.Protocol {
	case "udp":
		s.StartUDP()
	case "tcp":
		s.StartTCP()
	default:
		log.Fatal("bad protocol")
	}
}
