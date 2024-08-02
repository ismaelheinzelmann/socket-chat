package types

import (
	"net"
	"sync"
)

type User struct {
	Connection     *net.Conn
	Name           string
	ConnectionLock sync.RWMutex
}
