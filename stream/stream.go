package stream

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/JoshuaHenriques/proxy-server/listener"
)

type Stream struct {
	ClientIP   string
	ClientPort string
	ServerIP   string
	ServerPort string
	Protocol   string
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
	l, err := listener.New(s.Protocol, s.ClientPort)
	if err != nil {
		log.Fatal(err)
	}

	switch l := l.(type) {
	case *listener.UDPListener:
		l.Run()
		go s.handleUDPStream(l.Conn())
	case *listener.TCPListener:
		go l.Run()
		for lconn := range l.ConnChan() {
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

func (s *Stream) handleUDPStream(lconn *net.UDPConn) {
	defer lconn.Close()
	defer fmt.Printf("UDP Stream (%v) Finished/Closed\n------------------------------------\n", lconn)

	fmt.Printf("UDP Stream (%v) Running\n", lconn)

	clientRecv := make(chan Packet)

	go func() {
		for packet := range clientRecv {
			_, err := lconn.WriteToUDP(packet.Data, packet.Addr)
			if err != nil {
				fmt.Println(err)
			}
		}
	}()

	clientMap := make(map[string]chan Packet)
	buf := make([]byte, 65507)

	for {
		n, srcAddr, err := lconn.ReadFromUDP(buf)
		if err != nil {
			fmt.Printf("Connection Closed: %v\n", lconn)
			return
		}

		clientSend, ok := clientMap[srcAddr.String()]
		if ok {
			clientSend <- Packet{Data: buf[:n], Addr: srcAddr}
		} else {
			fmt.Printf("srcAddr: %s\n", srcAddr)
			sender := make(chan Packet)
			clientMap[srcAddr.String()] = sender

			go func() {
				var d net.Dialer
				conn, err := d.Dial(s.Protocol, fmt.Sprintf("%s:%s", s.ServerIP, s.ServerPort))
				if err != nil {
					fmt.Printf("failed to dial: %s", err)
				}
				dconn, ok := conn.(*net.UDPConn)
				if !ok {
					fmt.Printf("failed to dial: %s", err)
				}

				go func() {
					for {
						n, _, err := dconn.ReadFromUDP(buf)
						if err != nil {
							fmt.Printf("Connection Closed: %v\n", dconn)
							return
						}
						clientRecv <- Packet{Data: buf[:n], Addr: srcAddr}
					}
				}()

				for req := range sender {
					_, err = dconn.Write(req.Data)
					if err != nil {
						fmt.Println(err)
					}
				}
			}()
		}
	}
}
