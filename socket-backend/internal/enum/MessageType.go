package enum

type messageType struct {
	ErrorMessage    uint8
	OkMessage       uint8
	JoinMessage     uint8
	ChannelMessages uint8
	LeaveMessage    uint8
	WritingMessage  uint8
	CreateMessage   uint8
	SyncJoined      uint8
}

var MessageTypes = messageType{
	JoinMessage:     0,
	ChannelMessages: 1,
	LeaveMessage:    2,
	WritingMessage:  3,
	CreateMessage:   4,
	ErrorMessage:    5,
	OkMessage:       6,
	SyncJoined:      7,
}
