package types

import (
	"net"
	"sync"
)

type Channel struct {
	ID           uint8
	Members      map[*net.Conn]*User
	MembersLock  sync.RWMutex
	Name         string
	Messages     []Message
	MessagesLock sync.RWMutex
}
