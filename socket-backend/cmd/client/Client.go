package main

import (
	"encoding/json"
	"fmt"
	"net"
	"socket-backend/internal/enum"
	"socket-backend/internal/message/channel"
	"socket-backend/internal/message/common"
)

func main() {
	conn, err := net.Dial("tcp", ":8080")
	decoder := json.NewDecoder(conn)
	if err != nil {
		fmt.Println(err)
	}
	createMessage := channel.CreateChannelMessage{Name: "TEST"}
	cm, _ := json.Marshal(createMessage)
	message := common.PayloadMessage{MessageType: uint8(4), ChannelID: 0, Payload: &cm}
	m, _ := json.Marshal(message)
	_, _ = conn.Write(m)

	var p1 common.PayloadMessage
	_ = decoder.Decode(&p1)

	joinMessage := channel.JoinMessage{Name: "ISMAEL HEHEHEHEH"}
	cm, _ = json.Marshal(joinMessage)
	message = common.PayloadMessage{MessageType: enum.MessageTypes.JoinMessage, ChannelID: 1, Payload: &cm}
	m, _ = json.Marshal(message)
	_, _ = conn.Write(m)
	var p2 common.PayloadMessage
	_ = decoder.Decode(&p2)

}
