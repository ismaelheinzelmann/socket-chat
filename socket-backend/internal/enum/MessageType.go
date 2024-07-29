package enum

type messageType struct {
	JoinMessage     uint8
	ChannelMessages uint8
	LeaveMessage    uint8
	WritingMessage  uint8
}

var MessageTypes messageType = messageType{
	JoinMessage:     0,
	ChannelMessages: 1,
	LeaveMessage:    2,
	WritingMessage:  3,
}
