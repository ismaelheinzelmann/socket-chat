package commons

import (
	"net"
	"socket-backend/internal/service"
	"sync"
)

type Channel struct {
	Members      []net.Addr
	membersNames map[net.Addr]string
	name         string
	messages     []Message
	messagesLock sync.Mutex
	Handler      service.ChannelHandler
}
