package common

import (
	"net"
	"socket-backend/internal/message/common"
	"sync"
)

type Channel struct {
	ID           uint8
	members      []net.Addr
	membersNames map[net.Addr]string
	name         string
	messages     []common.Message
	messagesLock sync.Mutex
}
