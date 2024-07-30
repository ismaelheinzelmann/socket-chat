package handler

import (
	"net"
	"socket-backend/internal/commons"
	"socket-backend/internal/enum"
	"socket-backend/internal/message"
)

type ChannelHandler struct {
}

func (ch *ChannelHandler) Handle(channel *commons.Channel, connection *net.Conn, msg *message.Message) {
	switch msg.MessageType {
	case enum.MessageTypes.WritingMessage:
		{

		}
	case enum.MessageTypes.JoinMessage:
		{
			println("FUNCIONOU")
		}
	case enum.MessageTypes.LeaveMessage:
		{

		}
	case enum.MessageTypes.ChannelMessages:
		{

		}
	}
}
