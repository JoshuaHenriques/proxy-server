package stream

import (
	"fmt"
	"log"
	"net"

	"github.com/JoshuaHenriques/proxy-server/dialer"
	"github.com/JoshuaHenriques/proxy-server/listener"
)

type Stream struct {
	SrcIP    net.IP
	SrcPort  string
	DestIP   net.IP
	DestPort string
	Protocol string
}

func New(srcIPRaw, destIPRaw, srcPort, destPort, protocol string) *Stream {
	srcIP := net.ParseIP(srcIPRaw)
	destIP := net.ParseIP(destIPRaw)
	stream := &Stream{SrcIP: srcIP, DestIP: destIP, SrcPort: srcPort, DestPort: destPort, Protocol: protocol}
	return stream
}

func (s *Stream) Start() {
	srcChan := make(chan []byte)
	destChan := make(chan []byte)

	switch s.Protocol {
	case "udp":
		// go listener.UDPListener(bus, s.SrcPort)
		// go dialer.Dialer(bus, "udp", s.DestIP.String(), s.SrcPort)
	case "tcp":
		go listener.TCPListener(srcChan, destChan, s.SrcIP.String(), s.SrcPort)
		go dialer.Dialer(srcChan, destChan, "tcp", s.DestIP.String(), s.DestPort)
		// for elem := range bus {
		// 	fmt.Println("Data: ", string(elem))
		// }
	default:
		log.Fatal(fmt.Errorf("bad network protocol"))
	}
}
