package dialer

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"
)

func Dialer(bus chan []byte, protocol, ip, port string) {
	var d net.Dialer
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := d.DialContext(ctx, protocol, fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	for range bus {
		if _, err := conn.Write(<-bus); err != nil {
			log.Fatal(err)
		}
	}
}
