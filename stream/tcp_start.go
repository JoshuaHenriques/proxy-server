package stream

import (
	"bufio"
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

	lw := bufio.NewWriter(lconn)
	lr := bufio.NewReader(lconn)

	dw := bufio.NewWriter(dconn)
	dr := bufio.NewReader(dconn)

	// client <- server
	go func() {
		io.Copy(lw, dr)
		lw.Flush()
	}()
	// client -> server
	io.Copy(dw, lr)
	dw.Flush()
}
