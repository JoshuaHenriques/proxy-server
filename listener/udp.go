package listener

import (
	"fmt"
	"log"
	"net"
)

type UDPListener struct {
	conn *net.UDPConn
	port string
}

func (l *UDPListener) Run() {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%s", l.port))
	if err != nil {
		log.Fatal(err)
	}

	// unconnected socket, use WriteToUDP, can't use bufio
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatal(err)
	}
	// defer conn.Close()
	fmt.Printf("UDP Listener Conn: %v\n", conn)

	l.conn = conn
}

func (l *UDPListener) Conn() *net.UDPConn {
	return l.conn
}
