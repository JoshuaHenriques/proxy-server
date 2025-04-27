package main

import (
	"sync"

	"github.com/JoshuaHenriques/proxy-server/httpproxy"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

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

	go func() {
		httpProxy := httpproxy.NewProxy(
			"127.0.0.1",
			"192.168.2.18",
			"8080",
			"8080",
		)
		httpProxy.Start()
		wg.Done()
	}()

	go func() {
		httpsProxy := httpproxy.NewTLSProxy(
			"127.0.0.1",
			"192.168.2.18",
			"8080",
			"8080",
			"../cert.crt",
			"../cert.key",
		)
		httpsProxy.Start()
		wg.Done()
	}()

	wg.Wait()
}
