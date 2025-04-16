package dialer

import (
	"bufio"
	"fmt"
	"net"
)

type Dialer struct {
	Conn   net.Conn
	Reader *bufio.Reader
	Writer *bufio.Writer
}

func New(protocol, ip, port string) (*Dialer, error) {
	var d net.Dialer

	conn, err := d.Dial(protocol, fmt.Sprintf("%s:%s", ip, port))
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}
	fmt.Printf("Dialer Conn: %v\n", conn)

	w := bufio.NewWriter(conn)
	r := bufio.NewReader(conn)

	return &Dialer{Conn: conn, Reader: r, Writer: w}, nil
}
