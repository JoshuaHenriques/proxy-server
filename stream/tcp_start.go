package stream

import (
	"fmt"
	"io"
	"log"
	"net"
)

func (s *Stream) StartTCP() {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%s", s.ClientPort))
	if err != nil {
		log.Fatal(err)
	}

	tcpL, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	// l.SetDeadline(time.Now().Add(5 * time.Second))

	for {
		conn, err := tcpL.Accept()
		lconn := conn.(*net.TCPConn)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("TCP Listener Conn: %v\n", lconn)

		var d net.Dialer
		c, err := d.Dial(s.Protocol, fmt.Sprintf("%s:%s", s.ServerIP, s.ServerPort))
		if err != nil {
			log.Fatalf("failed to dial: %v\n", err)
		}

		dconn, ok := c.(*net.TCPConn)
		if !ok {
			log.Fatal("not tcpconn")
		}

		go s.handleTCPStream(lconn, dconn)
	}
}

func (s *Stream) handleTCPStream(lconn, dconn *net.TCPConn) {
	defer lconn.Close()
	defer dconn.Close()
	defer fmt.Printf("TCP Stream (%v) Finished/Closed\n------------------------------------\n", lconn)

	fmt.Printf("TCP Stream (%v) Running\n", lconn)

	// client <- server
	go io.Copy(lconn, dconn)
	// client -> server
	io.Copy(dconn, lconn)
}
