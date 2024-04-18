package model

type SendEventDto struct {
	ChatId string   `json:"chat_id"`
	Event  EventDto `json:"event"`
}

type EventDto struct {
	Type       string `json:"type"`
	Text       string `json:"text"`
	Visibility string `json:"visibility"`
}

func NewDefaultMessageEvent(chatId, body string) SendEventDto {
	event := EventDto{
		Type:       "message",
		Text:       body,
		Visibility: "all",
	}

	return SendEventDto{
		ChatId: chatId,
		Event:  event,
	}
}
