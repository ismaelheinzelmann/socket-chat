package message

type JoinMessage struct {
	ChannelId uint8  `json:"channelId"`
	Name      string `json:"name"`
}
