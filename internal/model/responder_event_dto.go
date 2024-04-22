package model

type ResponderEvent struct {
	ReplyType int
	ChatId    string
}

const (
	PLAIN_MESSAGE_REPLY = iota
	RICH_MESSAGE_REPLY
)

func NewPlainMessageResponderEvent(chatId string) ResponderEvent {
	return ResponderEvent{
		ReplyType: PLAIN_MESSAGE_REPLY,
		ChatId:    chatId,
	}
}

func NewRichMessageResponderEvent(chatId string) ResponderEvent {
	return ResponderEvent{
		ReplyType: RICH_MESSAGE_REPLY,
		ChatId:    chatId,
	}
}
