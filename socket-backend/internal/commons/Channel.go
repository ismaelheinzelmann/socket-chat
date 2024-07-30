package commons

import (
	"net"
	"sync"
)

type Channel struct {
	Members      []net.Addr
	membersNames map[net.Addr]string
	name         string
	messages     []Message
	messagesLock sync.Mutex
}
