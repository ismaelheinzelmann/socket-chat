package handler

import (
	"encoding/json"
	"net"
	"socket-backend/internal/Util"
	"socket-backend/internal/enum"
	"socket-backend/internal/message/channel"
	"socket-backend/internal/message/common"
	"socket-backend/internal/message/sync"
	"socket-backend/internal/types"
)

type ChannelHandler struct {
	channel types.Channel
}

func (ch *ChannelHandler) Handle(msg *common.PayloadMessage, conn *net.Conn) {
	var err *common.ErrorMessage
	switch msg.MessageType {
	case enum.MessageTypes.CreateMessage:

		err = ch.handleCreate(msg, conn)
	case enum.MessageTypes.JoinMessage:

		err = ch.handleJoin(msg, conn)
	case enum.MessageTypes.WritingMessage:
		err = ch.handleWriting(msg, conn)
	}
	if err != nil {
		payloadError := Util.ParseError(err)
		ch.sendError(payloadError, conn)
	}
}

//TODO Validar se o usuario esta no canal para mandar mensages

func (ch *ChannelHandler) handleCreate(msg *common.PayloadMessage, conn *net.Conn) *common.ErrorMessage {
	if ch.channel.ID != msg.ChannelID && ch.channel.ID != 0 {
		return &common.ErrorMessage{Error: "channel already exists"}
	}
	ch.channel.Members = make(map[*net.Conn]*types.User)
	var createMessage channel.CreateChannelMessage
	err := json.Unmarshal(*msg.Payload, &createMessage)
	if err != nil {
		return &common.ErrorMessage{Error: "error creating channel"}
	}
	ch.channel.ID = msg.ChannelID
	ch.channel.Name = createMessage.Name
	ch.sendOk("channel created successfully", conn)
	return nil
}

func (ch *ChannelHandler) handleJoin(msg *common.PayloadMessage, conn *net.Conn) *common.ErrorMessage {
	if ch.channel.ID != msg.ChannelID && ch.channel.ID != 0 {
		return &common.ErrorMessage{Error: "failed to join channel"}
	}
	var userJoined channel.JoinMessage
	err := json.Unmarshal(*msg.Payload, &userJoined)
	if err != nil {
		return &common.ErrorMessage{Error: "error joining channel"}
	}
	newMember := types.User{Name: userJoined.Name, Connection: conn}
	userJoinedMessage := sync.UserJoinedMessage{Name: userJoined.Name}
	userJoinedBytes, _ := json.Marshal(userJoinedMessage)
	payloadJoinedMessage := common.PayloadMessage{Payload: &userJoinedBytes, ChannelID: ch.channel.ID, MessageType: enum.MessageTypes.SyncJoined}
	payloadJoinedBytes, _ := json.Marshal(payloadJoinedMessage)
	ch.channel.MembersLock.RLock()
	for _, member := range ch.channel.Members {
		ch.sendMember(&payloadJoinedBytes, member)
	}
	ch.channel.Members[conn] = &newMember
	ch.sendMemberOk("channel joined successsfully", &newMember)
	ch.channel.MembersLock.Unlock()
	return nil
}

func (ch *ChannelHandler) handleWriting(msg *common.PayloadMessage, conn *net.Conn) *common.ErrorMessage {
	if ch.channel.ID != msg.ChannelID {
		return &common.ErrorMessage{Error: "failed to report status"}
	}
	var userWritingMessage sync.UserWritingMessage
	_ = json.Unmarshal(*msg.Payload, &userWritingMessage)
	payloadMessageBytes, _ := json.Marshal(userWritingMessage)
	ch.channel.MembersLock.RLock()
	defer ch.channel.MembersLock.RUnlock()
	for _, member := range ch.channel.Members {
		if member.Connection != conn {
			member.ConnectionLock.Lock()
			_, _ = (*member.Connection).Write(payloadMessageBytes)
			member.ConnectionLock.Unlock()
		}
	}
	return nil
}

func (ch *ChannelHandler) sendOk(message string, conn *net.Conn) {
	okMessage := common.OkMessage{Message: message}
	okBytes, _ := json.Marshal(okMessage)
	payloadMessage := common.PayloadMessage{Payload: &okBytes, MessageType: enum.MessageTypes.OkMessage, ChannelID: ch.channel.ID}
	payloadMessageBytes, _ := json.Marshal(payloadMessage)
	_, _ = (*conn).Write(payloadMessageBytes)
}

//	func (ch *ChannelHandler) sendMemberError(payloadError *common.PayloadMessage, user *types.User) {
//		(*user).ConnectionLock.Lock()
//		defer (*user).ConnectionLock.Unlock()
//		send
//	}
func (ch *ChannelHandler) sendMemberOk(message string, user *types.User) {
	(*user).ConnectionLock.Lock()
	defer (*user).ConnectionLock.Unlock()
	ch.sendOk(message, (*user).Connection)
}

func (ch *ChannelHandler) sendError(msg *common.PayloadMessage, conn *net.Conn) {
	payloadBytes, _ := json.Marshal(*msg)
	_, _ = (*conn).Write(payloadBytes)
}

func (ch *ChannelHandler) sendMember(messageBytes *[]byte, member *types.User) {
	(*member).ConnectionLock.Lock()
	defer (*member).ConnectionLock.Unlock()
	_, _ = (*(*member).Connection).Write(*messageBytes)
}

func (ch *ChannelHandler) verifyUserInChannel(conn *net.Conn) bool {
	_, ok := ch.channel.Members[conn]
	return ok
}
