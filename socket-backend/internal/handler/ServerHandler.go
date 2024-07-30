package handler

import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
	"socket-backend/internal/Util"
	"socket-backend/internal/commons"
)

type ServerHandler struct {
}

func (s *ServerHandler) Handle(conn *net.Conn) {
	defer (*conn).Close()
	reader := bufio.NewReader(*conn)
	for {
		// TODO: Implement error handling
		buffer := make([]byte, 4)
		_, err := io.ReadFull(reader, buffer)
		messageSize := binary.BigEndian.Uint32(buffer)
		messageBytes := make([]byte, messageSize)
		_, err = io.ReadFull(reader, messageBytes)
		msg, err := Util.Unmarshall(messageBytes)
		if err != nil {
			println(err)
		} else {
			println(msg.ChannelID)
		}
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
