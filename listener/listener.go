package listener

import (
	"fmt"
)

type Listener interface {
	Run()
}

func New(protocol, port string) (Listener, error) {
	if protocol != "udp" && protocol != "tcp" {
		return nil, fmt.Errorf("bad protocol")
	}

	switch protocol {
	case "udp":
		l := &UDPListener{port: port}
		return l, nil
	case "tcp":
		l := &TCPListener{port: port}
		return l, nil
	default:
		return nil, fmt.Errorf("bad protocol")
	}
}
