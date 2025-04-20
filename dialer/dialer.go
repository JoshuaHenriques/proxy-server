package dialer

import (
	"fmt"
	"net"
)

func New(protocol, ip, port string) (net.Conn, error) {
	var d net.Dialer

	conn, err := d.Dial(protocol, fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}
	fmt.Printf("Dialer Conn: %v\n", conn)

	return conn, nil
}
