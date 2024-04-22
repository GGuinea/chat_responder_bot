package model

type ResponderEvent struct {
	ActionToPerform int
	ChatId          string
	Action          string
}

const (
	PLAIN_MESSAGE_REPLY = iota
	RICH_MESSAGE_REPLY
	TRANSFER_TO_HUMAN_AGENT
)

var (
	ACTION_YES = "action_yes"
	ACTION_NO  = "action_no"
)

func NewPlainMessageResponderEvent(chatId string) ResponderEvent {
	return ResponderEvent{
		ActionToPerform: PLAIN_MESSAGE_REPLY,
		ChatId:          chatId,
	}
}

func NewRichMessageResponderEvent(chatId string) ResponderEvent {
	return ResponderEvent{
		ActionToPerform: RICH_MESSAGE_REPLY,
		ChatId:          chatId,
	}
}

func NewTransferResponderEvent(chatId string) ResponderEvent {
	return ResponderEvent{
		ActionToPerform: TRANSFER_TO_HUMAN_AGENT,
		ChatId:          chatId,
	}
}
