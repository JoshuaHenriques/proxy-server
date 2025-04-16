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
	Listener *listener.Listener
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
		fmt.Printf("Init UDP Listener\n")
		s.Listener = l

		s.run()
	case "tcp":
		l, err := listener.New(s.Protocol, s.SrcPort)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Init TCP Listener\n")
		s.Listener = l

	default:
		log.Fatal(fmt.Errorf("bad network protocol"))
	}

	s.run()
}

func (s *Stream) run() {
	s.Listener.Run()
	for conn := range s.Listener.ConnChan {
		d, err := dialer.New(s.Protocol, s.DestIP.String(), s.DestPort)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Init Dialer\n")

		switch s.Protocol {
		case "udp":
			go startUDPStream(conn, d)
		case "tcp":
			go startTCPStream(conn, d)
		}
	}
}

func startTCPStream(l *listener.Conn, d *dialer.Dialer) {
	defer l.Conn.Close()
	defer d.Conn.Close()
	defer fmt.Printf("TCP Stream (%v) Finished/Closed\n------------------------------------\n", l.Conn)

	fmt.Printf("TCP Stream (%v) Running\n", l.Conn)

	// client <- server
	go func() {
		io.Copy(l.Writer, d.Reader)
		l.Writer.Flush()
	}()
	// client -> server
	io.Copy(d.Writer, l.Reader)
	d.Writer.Flush()
}

func startUDPStream(l *listener.Conn, d *dialer.Dialer) {
	defer l.Conn.Close()
	defer d.Conn.Close()
	defer fmt.Printf("UDP Stream (%v) Finished/Closed\n------------------------------------\n", l.Conn)

	fmt.Printf("UDP Stream (%v) Running\n", l.Conn)

	// client <- server
	go func() {
		io.Copy(l.Writer, d.Reader)
		l.Writer.Flush()
	}()
	// client -> server
	io.Copy(d.Writer, l.Reader)
	d.Writer.Flush()
}
