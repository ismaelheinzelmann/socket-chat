package common

type Message struct {
	MessageType uint8   `json:"messageType"`
	ChannelID   uint8   `json:"channelId"`
	Payload     *[]byte `json:"payload"`
}
