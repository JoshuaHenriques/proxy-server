package proxy

import (
	"fmt"
	"net/http/httputil"
)

type Proxy struct {
	protocol, clientIP, clientPort, serverIP, serverPort, cert, key string
	reverseProxy                                                    *httputil.ReverseProxy
}

func New(protocol, clientIP, serverIP, clientPort, serverPort, cert, key string) *Proxy {
	proxy := &Proxy{
		protocol:   protocol,
		clientIP:   clientIP,
		clientPort: clientPort,
		serverIP:   serverIP,
		serverPort: serverPort,
		cert:       cert,
		key:        key,
	}

	return proxy
}

func (p *Proxy) Start() error {
	switch p.protocol {
	case "http":
		return p.StartHTTP()
	case "https":
		return p.StartHTTPS()
	default:
		return fmt.Errorf("bad protocol")
	}
}
