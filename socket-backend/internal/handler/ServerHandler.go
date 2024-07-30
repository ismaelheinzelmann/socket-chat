package handler

import (
	"bufio"
	"net"
	"socket-backend/internal/Util"
	"socket-backend/internal/commons"
)

type ServerHandler struct {
	channelHandler ChannelHandler
}

func (s *ServerHandler) Handle(conn *net.Conn) {
	defer (*conn).Close()
	reader := bufio.NewReader(*conn)
	for {
		_, e, eof := Util.ParseReceived(reader)
		if eof {
			break
		}
		if e != nil {
			_, err := (*conn).Write(Util.ParseSend(e))
			if err != nil {
				println(err.Error())
			}
			continue
		}

		//_, err := (*conn).Write(Util.ParseSend(msg))
		//if err != nil {
		//	println(err.Error())
		//}
	}

}

func (s *ServerHandler) getChannel(ch uint8) *commons.Channel {
	//for id, channel := range s.serverState.Channels {
	//	if id == ch {
	//		return channel
	//	}
	//}
	return nil
}
