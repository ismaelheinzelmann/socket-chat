package service

import (
	"net"
	"socket-backend/internal/enum"
	"socket-backend/internal/message"
)

type ChannelHandler struct {
	//channel commons.Channel
}

func (ch *ChannelHandler) Handle(connection net.Conn, msg message.Message) {
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
