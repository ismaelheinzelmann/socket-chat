package Util

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"io"
	"socket-backend/internal/message/common"
)

func Unmarshall(m []byte) (common.Message, error) {
	var msg common.Message
	err := json.Unmarshal(m, &msg)
	if err != nil {
		return msg, err
	}
	return msg, nil
}

func ParseReceived(reader *bufio.Reader) (*common.Message, *common.ErrorMessage, bool) {
	// TODO: implementar exception handler generalizado
	buffer := make([]byte, 4)
	_, err := io.ReadFull(reader, buffer)
	if err != nil {
		if err == io.EOF {
			return nil, nil, true
		}
		// TODO: Lidar com erro ao ler mensagem
	}
	messageSize := binary.BigEndian.Uint32(buffer)
	messageBytes := make([]byte, messageSize)
	_, err = io.ReadFull(reader, messageBytes)
	// TODO: Lidar com erro ao ler mensagem
	msg, err := Unmarshall(messageBytes)
	if err != nil {
		return &msg, &common.ErrorMessage{Error: err.Error()}, false
	}
	return &msg, nil, false
}

func ParseSend(msg any) []byte {
	m, _ := json.Marshal(msg)
	messageSizeArray := make([]byte, 4)
	binary.BigEndian.PutUint32(messageSizeArray, uint32(len(m)))
	m = append(messageSizeArray, m...)
	return m
}
