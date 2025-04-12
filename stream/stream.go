package stream

import (
	"fmt"
	"log"
	"net"

	"github.com/JoshuaHenriques/proxy-server/dialer"
	"github.com/JoshuaHenriques/proxy-server/listener"
)

type Stream struct {
	SrcPort  string
	DestIP   net.IP
	DestPort string
	Protocol string
}

func New(destIP string, srcPort string, destPort string, protocol string) *Stream {
	ip := net.ParseIP(destIP)
	stream := &Stream{DestIP: ip, SrcPort: srcPort, DestPort: destPort, Protocol: protocol}
	return stream
}

func (s *Stream) Start() {
	switch s.Protocol {
	case "udp":
		bus := make(chan []byte)
		go listener.UDPListener(bus, s.SrcPort)
		go dialer.Dialer(bus, "udp", s.DestIP.String(), s.SrcPort)
	case "tcp":
		bus := make(chan []byte)
		go listener.TCPListener(bus, s.SrcPort)
		// go dialer.Dialer(bus, "tcp", s.DestIP.String(), s.DestPort)
		for elem := range bus {
			fmt.Println("Data: ", string(elem))
		}
	default:
		// difference between panic and log.Fatal
		// panic(fmt.Errorf("bad network protocol"))
		log.Fatal(fmt.Errorf("bad network protocol"))
	}
}
