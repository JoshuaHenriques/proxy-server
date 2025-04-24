package stream

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/JoshuaHenriques/proxy-server/listener"
)

func (s *Stream) StartTCP() {
	l, err := listener.New(s.Protocol, s.ClientPort)
	if err != nil {
		log.Fatal(err)
	}
	tcpL := l.(*listener.TCPListener)
	go tcpL.Run()
	for lconn := range tcpL.ConnChan() {
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

type Packet struct {
	Data []byte
	Addr *net.UDPAddr
}
