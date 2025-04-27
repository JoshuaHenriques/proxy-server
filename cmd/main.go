package main

import (
	"log"
	"sync"

	"github.com/JoshuaHenriques/proxy-server/proxy"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)

	/*
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
	*/

	// go func() {
	// 	httpProxy := proxy.New(
	// 		"http",
	// 		"127.0.0.1",
	// 		"192.168.2.18",
	// 		"8080",
	// 		"8080",
	// 		"",
	// 		"",
	// 	)
	// 	if err := httpProxy.Start(); err != nil {
	// 		log.Printf("Proxy error: %v", err)
	// 	}
	// 	wg.Done()
	// }()

	go func() {
		httpsProxy := proxy.New(
			"https",
			"127.0.0.1",
			"192.168.2.18",
			"8080",
			"8080",
			"./cert.crt",
			"./cert.key",
		)
		if err := httpsProxy.Start(); err != nil {
			log.Fatal(err)
		}
		wg.Done()
	}()

	wg.Wait()
}
