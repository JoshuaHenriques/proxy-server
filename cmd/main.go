package main

import (
	"fmt"

	"github.com/JoshuaHenriques/proxy-server/stream"
)

func main() {
	stream := stream.New("127.0.0.1", "192.168.2.18", "7777", "7777", "tcp")
	fmt.Printf("stream: %+v\n", stream)
	stream.Start()

	select {}
}
