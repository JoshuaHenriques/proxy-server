package stream

import (
	"fmt"
	"io"
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
	switch s.Protocol {
	case "udp":
		l, err := listener.New(s.Protocol, s.SrcPort)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Started UDP Listener\n")

		d, err := dialer.New(s.Protocol, s.DestIP.String(), s.DestPort)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Started Dialer\n")

		go startUDPStream(l, d)
	case "tcp":
		l, err := listener.New(s.Protocol, s.SrcPort)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Started TCP Listener\n")

		d, err := dialer.New(s.Protocol, s.DestIP.String(), s.DestPort)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Started Dialer\n")

		go startTCPStream(l, d)
	default:
		log.Fatal(fmt.Errorf("bad network protocol"))
	}
}

func startTCPStream(l *listener.Listener, d *dialer.Dialer) {
	defer l.Conn.Close()
	defer d.Conn.Close()

	// l.Run()

	fmt.Printf("TCP Stream Running\n")
	go func() {
		io.Copy(d.Writer, l.Reader)
		d.Writer.Flush()
	}()
	io.Copy(l.Writer, d.Reader)
	l.Writer.Flush()
}

func startUDPStream(l *listener.Listener, d *dialer.Dialer) {
	defer l.Conn.Close()
	defer d.Conn.Close()

	fmt.Printf("UDP Stream Running\n")
	go func() {
		io.Copy(d.Writer, l.Reader)
		d.Writer.Flush()
	}()
	io.Copy(l.Writer, d.Reader)
	l.Writer.Flush()
}
