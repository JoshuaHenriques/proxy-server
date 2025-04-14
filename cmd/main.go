package main

import (
	"github.com/JoshuaHenriques/proxy-server/stream"
)

func main() {
	stream := stream.New("127.0.0.1", "192.168.2.18", "7777", "7777", "tcp")
	stream.Start()

	select {}
}
