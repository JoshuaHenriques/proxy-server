package httpproxy

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type TLSProxy struct {
	clientIP, clientPort, serverIP, serverPort, cert, key string
	reverseProxy                                          *httputil.ReverseProxy
}

func NewTLSProxy(clientIP, serverIP, clientPort, serverPort, cert, key string) *TLSProxy {
	proxy := &TLSProxy{
		clientIP:   clientIP,
		clientPort: clientPort,
		serverIP:   serverIP,
		serverPort: serverPort,
		cert:       cert,
		key:        key,
	}

	return proxy
}

func (p *TLSProxy) Start() error {
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

	tlsConfig := &tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.CurveP256, tls.X25519},
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", p.handleProxy)

	server := &http.Server{
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	srvURL, err := url.Parse(
		fmt.Sprintf("http://%s:%s", p.serverIP, p.serverPort))
	if err != nil {
		return fmt.Errorf("error parsing url, err: %s\n", err)
	}

	p.reverseProxy = httputil.NewSingleHostReverseProxy(srvURL)
	p.reverseProxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("proxy error: %v", err)
		w.WriteHeader(http.StatusBadGateway)
	}

	log.Printf("HTTPS Proxy Running on :%s\n", p.clientPort)
	return server.ServeTLS(tcpL, p.cert, p.key)
}

func (p *TLSProxy) handleProxy(w http.ResponseWriter, req *http.Request) {
	p.reverseProxy.ServeHTTP(w, req)
}
