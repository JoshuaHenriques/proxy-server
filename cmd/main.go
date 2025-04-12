package main

import (
	"fmt"

	"github.com/JoshuaHenriques/proxy-server/stream"
)

func main() {
	stream := stream.New("127.0.0.1", "8888", "7777", "tcp")
	fmt.Printf("stream: %+v\n", stream)
	stream.Start()
}
