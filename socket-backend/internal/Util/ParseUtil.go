package Util

import (
	"encoding/json"
	"socket-backend/internal/message"
)

func Unmarshall(m []byte) (message.Message, error) {
	var msg message.Message
	err := json.Unmarshal(m, &msg)
	if err != nil {
		return msg, err
	}
	return msg, nil
}
