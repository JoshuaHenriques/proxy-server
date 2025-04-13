package dialer

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func Dialer(srcChan, destChan chan []byte, protocol, ip, port string) {
	var d net.Dialer

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	fmt.Printf("ip: %s, port: %s\n", ip, port)
	conn, err := d.DialContext(ctx, protocol, fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	go outboundHandler(srcChan, conn)
	// go inboundHandler(destChan, conn)
}

func outboundHandler(srcChan chan []byte, conn net.Conn) {
	for range srcChan {
		if _, err := conn.Write(<-srcChan); err != nil {
			log.Fatal(err)
		}
	}
}

func inboundHandler(destChan chan []byte, conn net.Conn) {
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
		destChan <- bytes
	}
}
