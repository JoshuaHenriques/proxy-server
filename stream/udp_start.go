package stream

import (
	"fmt"
	"log"
	"net"
)

type Packet struct {
	Data []byte
	Addr *net.UDPAddr
}

func (s *Stream) StartUDP() {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%s", s.ClientPort))
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("UDP Listener Conn: %v\n", conn)

	go s.handleUDPStream(conn)
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
