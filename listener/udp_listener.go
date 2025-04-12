package listener

import (
	"fmt"
	"log"
	"net"
)

func UDPListener(bus chan []byte, port string) {
	l, err := net.ListenPacket("udp", port)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	fmt.Println("Init UDP Listener")

	// safely above MTU of 1500 bytes
	buff := make([]byte, 4096)
	for {
		n, _, err := l.ReadFrom(buff)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Bytes read: %d\n", n)

		// todo: udp handler that streams each packet to destination ip:port
	}
}
