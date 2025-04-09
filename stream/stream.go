package stream

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"
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
	// todo: initialize channel

	if s.Protocol == "udp" {
		// spawn udpListener goroutine
		// go s.initUDPListener()
		s.initUDPListener()
		// listener.UDPListener()
	}

	if s.Protocol == "tcp" {
		// spawn tcpListener goroutine
		// go s.initTCPListener()
		s.initTCPListener()
		// listener.TCPListener()
	}
	// go s.initListener()
	// go s.initDialer()
}

func (s *Stream) initTCPListener() {
	l, err := net.Listen("tcp", s.SrcPort)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	fmt.Println("Init TCP Listener")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		io.Copy(conn, conn)
		fmt.Printf("conn: %+v\n", conn)
		// go s.tcpConnectionHandler(conn)
	}
}

func (s *Stream) tcpConnectionHandler(c net.Conn) {
	// todo: open channel and write to it

	io.Copy(c, c)
	c.Close()
}

func (s *Stream) initUDPListener() {
	l, err := net.ListenPacket("udp", s.SrcPort)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	fmt.Println("Init UDP Listener")

	// safely above MTU of 1500 bytes
	buff := make([]byte, 4096)
	for {
		n, _, err := l.ReadFrom(buff)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Bytes read: %d", n)

		// todo: udp handler that streams each packet to destination ip:port
	}
}

func (s *Stream) initDialer() {
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := d.DialContext(ctx, s.DestPort, fmt.Sprintf("%s:%s", s.DestIP.String(), s.DestPort))
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	fmt.Printf("Init Dialer")
	// todo: open channel and read it then pass it to conn.Write

	if _, err := conn.Write([]byte("Hello, World!")); err != nil {
		log.Fatal(err)
	}
}
