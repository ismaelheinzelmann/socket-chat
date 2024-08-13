package types

import "net"

type Message struct {
	Origin *net.Conn
	Body   string
}
