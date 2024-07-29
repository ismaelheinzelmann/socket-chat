package Socket

import (
	"bufio"
	"fmt"
	"net"
	"socket-backend/internal/Util"
	"socket-backend/internal/commons"
	"socket-backend/internal/enum"
	"sync"
)

type Server struct {
	serverState commons.State
	stateLock   sync.Mutex
	port        string
}

func NewServer(port string) *Server {
	return &Server{serverState: commons.State{}, port: port, stateLock: sync.Mutex{}}
}

func (s *Server) Run() {
	fmt.Println("Listening on port " + s.port)
	ln, _ := net.Listen("tcp", ":"+s.port)

	for {
		//var m []byte
		conn, _ := ln.Accept()
		reader := bufio.NewReader(conn)
		data, err := reader.ReadBytes('\n')
		s.stateLock.Lock()
		msg, err := Util.Unmarshall(data)
		if err != nil {
			println(err)
		}
		if msg.MessageType == enum.MessageTypes.JoinMessage {
			if err != nil {
				fmt.Println(err)
			}
			channel := s.getChannel(msg.ChannelID)
			if channel != nil {
				s.serverState.Channels[msg.ChannelID].Handler.Handle(conn, msg)
			}
		}
		// TODO: Lidar com outros
		s.stateLock.Unlock()
	}
}

func (s *Server) getChannel(ch uint8) *commons.Channel {
	for id, channel := range s.serverState.Channels {
		if id == ch {
			return channel
		}
	}
	return nil
}
