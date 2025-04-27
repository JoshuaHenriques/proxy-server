package main

import (
	"sync"

	"github.com/JoshuaHenriques/proxy-server/stream"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		tcpStream := stream.New(
			"127.0.0.1",
			"192.168.2.18",
			"7777",
			"7777",
			"tcp")
		tcpStream.Start()
		wg.Done()
	}()

	go func() {
		udpStream := stream.New(
			"127.0.0.1",
			"192.168.2.18",
			"7777",
			"7777",
			"udp")
		udpStream.Start()
		wg.Done()
	}()

	// httpProxy := httpproxy.New(
	// 	"127.0.0.1",
	// 	"192.168.2.18",
	// 	"7777",
	// 	"7777",
	// )
	// go httpProxy.Start()

	// select {}
	wg.Wait()
}
