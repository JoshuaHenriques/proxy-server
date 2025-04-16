package listener

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type Listener struct {
	ConnChan chan *Conn
	Protocol string
	Port     string
}

type Conn struct {
	Conn   net.Conn
	Reader *bufio.Reader
	Writer *bufio.Writer
}

func New(protocol, port string) (*Listener, error) {
	listener := &Listener{Port: port, Protocol: protocol}
	listener.ConnChan = make(chan *Conn)

	if protocol != "udp" && protocol != "tcp" {
		return nil, fmt.Errorf("bad protocol")
	}
	return listener, nil
}

func (l *Listener) Run() {
	switch l.Protocol {
	case "udp":
		go l.startUDPListener()
	case "tcp":
		go l.startTCPListener()
	}
}

func (l *Listener) startTCPListener() {
	addr, err := net.ResolveTCPAddr(l.Protocol, fmt.Sprintf(":%s", l.Port))
	if err != nil {
		log.Fatal(err)
	}

	tcpL, err := net.ListenTCP(l.Protocol, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer tcpL.Close()
	// l.SetDeadline(time.Now().Add(5 * time.Second))

	for {
		conn, err := tcpL.Accept()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("TCP Listener Conn: %v\n", conn)

		w := bufio.NewWriter(conn)
		r := bufio.NewReader(conn)

		l.ConnChan <- &Conn{Conn: conn, Reader: r, Writer: w}

		// l.Conn = conn
		// l.Reader = r
		// l.Writer = w
	}
}

func (l *Listener) startUDPListener() {
	addr, err := net.ResolveUDPAddr(l.Protocol, fmt.Sprintf(":%s", l.Port))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP(l.Protocol, addr)
	if err != nil {
		log.Fatal(err)
	}
	// defer conn.Close()
	fmt.Printf("UDP Listener Conn: %v\n", conn)

	w := bufio.NewWriter(conn)
	r := bufio.NewReader(conn)

	l.ConnChan <- &Conn{Conn: conn, Reader: r, Writer: w}
	close(l.ConnChan)

	// l.Conn = conn
	// l.Reader = r
	// l.Writer = w
}
