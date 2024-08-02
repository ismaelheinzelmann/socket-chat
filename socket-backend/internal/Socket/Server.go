package Socket

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"socket-backend/internal/enum"
	"socket-backend/internal/handler"
	"socket-backend/internal/message/common"
	"sync"
)

type Server struct {
	port      string
	handlers  map[uint8]*handler.ChannelHandler
	stateLock sync.RWMutex
}

func NewServer(port string) *Server {
	return &Server{port: port}
}

func (s *Server) Run() {
	s.handlers = make(map[uint8]*handler.ChannelHandler)
	fmt.Println("Listening on port " + s.port)
	ln, _ := net.Listen("tcp", ":"+s.port)
	for {
		conn, _ := ln.Accept()
		go s.Handle(&conn)
	}
}

func (s *Server) Handle(conn *net.Conn) {
	defer (*conn).Close()
	decoder := json.NewDecoder(*conn)
	for {
		// TODO refatorar verbosidade dos erros

		var payloadMessage common.PayloadMessage
		err := decoder.Decode(&payloadMessage)
		if err != nil {
			if err == io.EOF {
				break
			}
		}
		channelHandler, _ := s.getHandler(&payloadMessage)
		channelHandler.Handle(&payloadMessage, conn)
	}
}

func (s *Server) getHandler(msg *common.PayloadMessage) (*handler.ChannelHandler, *common.ErrorMessage) {
	if msg.MessageType != enum.MessageTypes.CreateMessage && msg.ChannelID == 0 {
		return nil, &common.ErrorMessage{Error: "invalid channel"}
	}
	if msg.MessageType == enum.MessageTypes.CreateMessage && msg.ChannelID != 0 {
		return nil, &common.ErrorMessage{Error: "invalid channel"}
	}
	if msg.MessageType == enum.MessageTypes.CreateMessage && msg.ChannelID == 0 {
		id, newHandler := s.newChannel()
		msg.ChannelID = id
		return newHandler, nil
	}
	s.stateLock.RLock()
	channelHandler, ok := s.handlers[msg.ChannelID]
	if !ok {
		return nil, &common.ErrorMessage{Error: "invalid channel"}
	}
	s.stateLock.RUnlock()
	return channelHandler, nil
}

// TODO quando ultimo membro sair delete o canal ?

func (s *Server) newChannel() (uint8, *handler.ChannelHandler) {
	s.stateLock.Lock()
	defer s.stateLock.Unlock()
	ch := handler.ChannelHandler{}
	id := s.getNextId()
	s.handlers[id] = &ch
	return id, &ch
}

func (s *Server) getNextId() uint8 {
	i := 1
	for {
		_, ok := s.handlers[uint8(i)]
		if !ok {
			return uint8(i)
		}
		i++
	}

}
