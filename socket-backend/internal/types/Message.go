package common

import "net"

type Message struct {
	origin net.Addr
	body   string
}
