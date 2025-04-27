package proxy

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func (p *Proxy) StartHTTP() error {
	addr, err := net.ResolveTCPAddr(
		"tcp",
		fmt.Sprintf(":%s", p.clientPort))
	if err != nil {
		return fmt.Errorf("invalid server URL: %s", err)
	}

	tcpL, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return fmt.Errorf("error with tcp listener: %s", err)
	}
	defer tcpL.Close()

	srvURL, err := url.Parse(
		fmt.Sprintf("http://%s:%s", p.serverIP, p.serverPort))
	if err != nil {
		fmt.Printf("error parsing url, err: %s\n", err)
	}

	p.reverseProxy = httputil.NewSingleHostReverseProxy(srvURL)
	p.reverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("proxy error: %v", err)
		w.WriteHeader(http.StatusBadGateway)
	}

	log.Printf("HTTP Proxy Running on :%s\n", p.clientPort)
	mux := http.NewServeMux()
	mux.HandleFunc("/", p.handleProxy)
	return http.Serve(tcpL, mux)
}

func (p *Proxy) handleProxy(w http.ResponseWriter, req *http.Request) {
	p.reverseProxy.ServeHTTP(w, req)
}
