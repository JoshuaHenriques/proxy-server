package stream

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/JoshuaHenriques/proxy-server/dialer"
	"github.com/JoshuaHenriques/proxy-server/listener"
)

type Stream struct {
	ClientIP   string
	ClientPort string
	ServerIP   string
	ServerPort string
	Protocol   string
	Listener   *listener.Listener
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
		l, err := listener.New(s.Protocol, s.ClientPort)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Init UDP Listener\n")
		s.Listener = l

	case "tcp":
		l, err := listener.New(s.Protocol, s.ServerPort)
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
		d, err := dialer.New(s.Protocol, s.ServerIP, s.ServerPort)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Init Dialer\n")

		switch s.Protocol {
		case "udp":
			lconn, ok := conn.(*net.UDPConn)
			if !ok {
				log.Fatal("not udpconn")
			}

			dconn, ok := d.(*net.UDPConn)
			if !ok {
				log.Fatal("not udpconn")
			}

			go s.startUDPStream(lconn, dconn)
			close(s.Listener.ConnChan)
		case "tcp":
			lconn, ok := conn.(*net.TCPConn)
			if !ok {
				log.Fatal("not tcpconn")
			}

			dconn, ok := d.(*net.TCPConn)
			if !ok {
				log.Fatal("not tcpconn")
			}

			go s.startTCPStream(lconn, dconn)
		}
	}
}

func (s *Stream) startTCPStream(lconn, dconn *net.TCPConn) {
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

func (s *Stream) startUDPStream(lconn, dconn *net.UDPConn) {
	defer lconn.Close()
	defer dconn.Close()
	defer fmt.Printf("UDP Stream (%v) Finished/Closed\n------------------------------------\n", lconn)

	fmt.Printf("UDP Stream (%v) Running\n", lconn)

	var clientAddr *net.UDPAddr
	buf := make([]byte, 65507)

	// multi client
	// keep track of all clientaddr and send

	// client <- server
	go func() {
		for {
			n, _, err := dconn.ReadFromUDP(buf)
			if err != nil {
				fmt.Printf("Connection Closed: %v\n", dconn)
				return
			}
			// fmt.Printf("%s\n", string(buf[:n]))
			_, err = lconn.WriteToUDP(buf[:n], clientAddr)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	// client -> server
	for {
		n, addr, err := lconn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("Connection Closed: %v\n", lconn)
			return
		}

		clientAddr = addr

		_, err = dconn.Write(buf[:n])
		if err != nil {
			fmt.Println(err)
		}
	}
}
