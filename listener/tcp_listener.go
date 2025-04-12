package listener

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func TCPListener(bus chan []byte, port string) {
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%s", port))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	fmt.Println("Init TCP Listener")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
			close(bus)
			break // return
		}
		fmt.Printf("new conn: %+v\n", conn)
		go tcpHandler(bus, conn)
	}
}

func tcpHandler(bus chan []byte, conn net.Conn) {
	r := bufio.NewReader(conn)
	for {
		bytes := make([]byte, 4096)
		_, err := r.Read(bytes)
		switch err {
		case nil:
		case io.EOF:
			fmt.Println("EOF", err)
			conn.Close()
			return
		default:
			fmt.Println("ERROR", err)
			conn.Close()
			return
		}
		bus <- bytes
	}
}
