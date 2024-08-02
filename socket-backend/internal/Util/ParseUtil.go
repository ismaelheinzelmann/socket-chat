package Util

import (
	"encoding/json"
	"socket-backend/internal/enum"
	"socket-backend/internal/message/common"
)

func Unmarshall(m []byte) (common.PayloadMessage, error) {
	var msg common.PayloadMessage
	err := json.Unmarshal(m, &msg)
	if err != nil {
		return msg, err
	}
	return msg, nil
}

func ParseError(err *common.ErrorMessage) *common.PayloadMessage {
	errBytes, _ := json.Marshal(*err)
	return &common.PayloadMessage{Payload: &errBytes, ChannelID: 0, MessageType: enum.MessageTypes.ErrorMessage}
}
