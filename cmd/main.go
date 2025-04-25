package main

import (
	"github.com/JoshuaHenriques/proxy-server/stream"
)

func main() {
	tcpStream := stream.New(
		"127.0.0.1",
		"192.168.2.18",
		"7777",
		"7777",
		"tcp")
	go tcpStream.Start()

	udpStream := stream.New(
		"127.0.0.1",
		"192.168.2.18",
		"7777",
		"7777",
		"udp")
	go udpStream.Start()

	// httpProxy := httpproxy.New(
	// 	"127.0.0.1",
	// 	"192.168.2.18",
	// 	"7777",
	// 	"7777",
	// )
	// go httpProxy.Start()

	select {}
}
