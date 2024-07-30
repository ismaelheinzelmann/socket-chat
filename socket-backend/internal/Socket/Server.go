package Socket

import (
	"fmt"
	"net"
	"socket-backend/internal/commons"
	"socket-backend/internal/handler"
	"sync"
)

type Server struct {
	serverState    commons.State
	stateLock      sync.Mutex
	port           string
	channelHandler handler.ChannelHandler
	serverHandler  handler.ServerHandler
}

func NewServer(port string) *Server {
	return &Server{serverState: commons.State{}, port: port, stateLock: sync.Mutex{}}
}

func (s *Server) Run() {
	fmt.Println("Listening on port " + s.port)
	ln, _ := net.Listen("tcp", ":"+s.port)

	for {
		conn, _ := ln.Accept()
		go s.serverHandler.Handle(&conn)
	}
}

func (s *Server) receiveConnection(conn net.Conn) {
}
