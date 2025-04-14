package listener

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func TCPListener(srcChan, destChan chan []byte, ip, port string) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	fmt.Println("Init TCP Listener")
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("new conn: %+v\n", conn)
		go outboundHandler(destChan, conn)
		go inboundHandler(srcChan, conn)
	}
}

func outboundHandler(destChan chan []byte, conn net.Conn) {
	defer conn.Close()

	w := bufio.NewWriter(conn)
	for bytes := range destChan {
		if _, err := w.Write(bytes); err != nil {
			log.Fatal(err)
			return
		}
	}
}

func inboundHandler(srcChan chan []byte, conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	for {
		bytes := make([]byte, 4096)
		n, err := r.Read(bytes)

		switch err {
		case nil:
			fmt.Printf("bytes read: %d\n", n)
			srcChan <- bytes[:n]
		case io.EOF:
			fmt.Println("EOF", err)
			conn.Close()
			return
		default:
			fmt.Println("ERROR", err)
			conn.Close()
			return
		}
	}
}
