package listener

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

func TCPListener(srcChan, destChan chan []byte, ip, port string) {
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", ip, port))
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
		// go outboundHandler(destChan, conn)
		go inboundHandler(srcChan, conn)
	}
}

func outboundHandler(destChan chan []byte, conn net.Conn) {
	for range destChan {
		if _, err := conn.Write(<-destChan); err != nil {
			log.Fatal(err)
		}
	}
}

func inboundHandler(srcChan chan []byte, conn net.Conn) {
	r := bufio.NewReader(conn)
	for {
		bytes := make([]byte, 4096)
		_, err := r.Read(bytes)
		switch err {
		case nil:
		case io.EOF:
			fmt.Println("EOF", err)
			// conn.Close()
			// return
		default:
			fmt.Println("ERROR", err)
			conn.Close()
			return
		}
		srcChan <- bytes
	}
}
