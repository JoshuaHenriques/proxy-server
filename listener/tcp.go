package listener

import (
	"fmt"
	"log"
	"net"
)

type TCPListener struct {
	connChan chan *net.TCPConn
	port     string
}

func (l *TCPListener) Run() {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%s", l.port))
	if err != nil {
		log.Fatal(err)
	}

	tcpL, err := net.ListenTCP("tcp", addr)
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

		l.connChan <- conn.(*net.TCPConn)
	}
}

func (l *TCPListener) ConnChan() <-chan *net.TCPConn {
	return l.connChan
}
