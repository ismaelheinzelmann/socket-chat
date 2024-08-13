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

func (ch *ChannelHandler) GetInfo() (string, uint16) {
	ch.channel.MembersLock.RLock()
	defer ch.channel.MembersLock.RUnlock()
	var members uint16 = 0
	for _, _ = range ch.channel.Members {
		members += 1
	}
	return ch.channel.Name, members
}

func (ch *ChannelHandler) Handle(msg *common.PayloadMessage, conn *net.Conn) bool {
	var err *common.ErrorMessage
	var keep bool
	switch msg.MessageType {
	case enum.MessageTypes.CreateMessage:

		keep, err = ch.handleCreate(msg, conn)
	case enum.MessageTypes.JoinMessage:

		keep, err = ch.handleJoin(msg, conn)
	case enum.MessageTypes.WritingMessage:
		keep, err = ch.handleWriting(msg, conn)
	case enum.MessageTypes.SendMessage:
		keep, err = ch.handleSendMessage(msg, conn)
	case enum.MessageTypes.LeaveMessage:
		keep, err = ch.handleLeave(msg, conn)
	}

	if err != nil {
		payloadError := Util.ParseError(err)
		ch.sendError(payloadError, conn)
	}
	return keep
}

//TODO Validar se o usuario esta no canal para mandar mensages

func (ch *ChannelHandler) handleCreate(msg *common.PayloadMessage, conn *net.Conn) (bool, *common.ErrorMessage) {
	if ch.channel.ID != msg.ChannelID && ch.channel.ID != 0 {
		return true, &common.ErrorMessage{Error: "channel already exists"}
	}
	ch.channel.Members = make(map[*net.Conn]*types.User)
	var createMessage channel.CreateChannelMessage
	err := json.Unmarshal(*msg.Payload, &createMessage)
	if err != nil {
		return true, &common.ErrorMessage{Error: "error creating channel"}
	}
	ch.channel.ID = msg.ChannelID
	ch.channel.Name = createMessage.Name
	ch.sendOk("channel created successfully", conn)
	return true, nil
}

func (ch *ChannelHandler) handleJoin(msg *common.PayloadMessage, conn *net.Conn) (bool, *common.ErrorMessage) {
	if ch.channel.ID != msg.ChannelID && ch.channel.ID != 0 {
		return true, &common.ErrorMessage{Error: "failed to join channel"}
	}
	var userJoined channel.JoinMessage
	err := json.Unmarshal(*msg.Payload, &userJoined)
	if err != nil {
		return true, &common.ErrorMessage{Error: "error joining channel"}
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
	ch.channel.MembersLock.RUnlock()
	return true, nil
}

func (ch *ChannelHandler) handleLeave(msg *common.PayloadMessage, conn *net.Conn) (bool, *common.ErrorMessage) {
	if ch.verifyUserInChannel(conn) {
		ch.channel.MembersLock.RLock()
		defer ch.channel.MembersLock.RUnlock()
		syncLeaveMessage := sync.UserLeaveMessage{Name: ch.channel.Members[conn].Name}
		syncLeaveBytes, _ := json.Marshal(syncLeaveMessage)
		payloadMessage := common.PayloadMessage{MessageType: enum.MessageTypes.SyncLeaveMessage, ChannelID: ch.channel.ID, Payload: &syncLeaveBytes}
		payloadMessageBytes, _ := json.Marshal(payloadMessage)
		for _, member := range ch.channel.Members {
			if member.Connection != conn {
				ch.sendMember(&payloadMessageBytes, member)
			}
		}
	}
	return false, nil
}

func (ch *ChannelHandler) handleWriting(msg *common.PayloadMessage, conn *net.Conn) (bool, *common.ErrorMessage) {
	if ch.verifyUserInChannel(conn) {
		if ch.channel.ID != msg.ChannelID {
			return true, &common.ErrorMessage{Error: "failed to report status"}
		}
		var userWritingMessage sync.UserWritingMessage
		_ = json.Unmarshal(*msg.Payload, &userWritingMessage)
		payloadMessageBytes, _ := json.Marshal(userWritingMessage)
		ch.channel.MembersLock.RLock()
		userWritingMessage.UserWriting = ch.channel.Members[conn].Name
		writingBytes, _ := json.Marshal(userWritingMessage)
		msg.Payload = &writingBytes
		defer ch.channel.MembersLock.RUnlock()
		for _, member := range ch.channel.Members {
			if member.Connection != conn {
				ch.sendMember(&payloadMessageBytes, member)
			}
		}
	}
	return true, nil
}

func (ch *ChannelHandler) handleSendMessage(msg *common.PayloadMessage, conn *net.Conn) (bool, *common.ErrorMessage) {
	if ch.verifyUserInChannel(conn) {
		if ch.channel.ID != msg.ChannelID {
			return true, &common.ErrorMessage{Error: "failed to report status"}
		}
		var sentMessage channel.MessageSendMessage
		_ = json.Unmarshal(*msg.Payload, &sentMessage)
		ch.channel.MembersLock.RLock()
		ch.channel.MessagesLock.RLock()
		defer ch.channel.MembersLock.RUnlock()
		defer ch.channel.MessagesLock.RUnlock()

		ch.channel.Messages = append(ch.channel.Messages, types.Message{Origin: conn, Body: sentMessage.Message})
		syncSendMessage := sync.SendMessage{Name: ch.channel.Members[conn].Name, Message: sentMessage.Message}
		syncSendMessageBytes, _ := json.Marshal(syncSendMessage)
		payloadMessage := common.PayloadMessage{MessageType: enum.MessageTypes.SyncSendMessage, ChannelID: ch.channel.ID, Payload: &syncSendMessageBytes}
		payloadMessageBytes, _ := json.Marshal(payloadMessage)
		for _, member := range ch.channel.Members {
			if member.Connection != conn {
				ch.sendMember(&payloadMessageBytes, member)
			}
		}
	}
	return true, nil
}

//TODO: Enum of messages

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

func (ch *ChannelHandler) DisconnectUser(conn *net.Conn) {
	ch.channel.MembersLock.Lock()
	defer ch.channel.MembersLock.Unlock()
	_, ok := ch.channel.Members[conn]
	if ok {
		delete(ch.channel.Members, conn)
	}
}
