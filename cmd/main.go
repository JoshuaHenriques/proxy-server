package main

import (
	"fmt"

	"github.com/JoshuaHenriques/proxy-server/stream"
)

func main() {
	stream := stream.New("localhost", ":8888", ":7777", "tcp")
	fmt.Printf("stream: %+v\n", stream)
	stream.Start()
}
