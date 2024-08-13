package enum

type messageType struct {
	ErrorMessage        uint8
	OkMessage           uint8
	JoinMessage         uint8
	ChannelMessages     uint8
	LeaveMessage        uint8
	WritingMessage      uint8
	CreateMessage       uint8
	SyncJoined          uint8
	ListChannelMessages uint8
	SendMessage         uint8
	SyncSendMessage     uint8
	SyncLeaveMessage    uint8
}

var MessageTypes = messageType{
	JoinMessage:         0,
	SyncJoined:          7,
	ChannelMessages:     1,
	LeaveMessage:        2,
	SyncLeaveMessage:    11,
	WritingMessage:      3,
	CreateMessage:       4,
	ErrorMessage:        5,
	OkMessage:           6,
	ListChannelMessages: 8,
	SendMessage:         9,
	SyncSendMessage:     10,
}

//List
//Create
//Join
//Leave
//Send message
//Receive message
//Send writing synchronization
//Receive writing synchronization
//Receive user leaving
