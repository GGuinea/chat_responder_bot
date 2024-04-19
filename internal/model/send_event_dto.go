package model

type EventInterface interface {
}

type SendEventDto struct {
	ChatId string         `json:"chat_id"`
	Event  EventInterface `json:"event"`
}

type MessageEventDto struct {
	Type       string `json:"type"`
	Text       string `json:"text"`
	Visibility string `json:"visibility"`
}

type RichMessageEventDto struct {
	Type       string               `json:"type"`
	TemplateId string               `json:"template_id"`
	Visibility string               `json:"visibility"`
	Elements   []RichMessageElement `json:"elements"`
}

type RichMessageElement struct {
	Title    string              `json:"title"`
	Subtitle string              `json:"subtitle"`
	Buttons  []RichMessageButton `json:"buttons"`
}

type RichMessageButton struct {
	Type       string   `json:"type"`
	Text       string   `json:"text"`
	PostbackId string   `json:"postback_id"`
	Value      string   `json:"value"`
	UserIDS    []string `json:"user_ids"`
}

func NewDefaultMessageEvent(chatId, body string) SendEventDto {
	event := MessageEventDto{
		Type:       "message",
		Text:       body,
		Visibility: "all",
	}

	return SendEventDto{
		ChatId: chatId,
		Event:  event,
	}
}

func NewRichCardMessageEvent(chatId string) SendEventDto {
	richMessage := RichMessageEventDto{
		Type:       "rich_message",
		TemplateId: "cards",
		Visibility: "all",
		Elements: []RichMessageElement{
			{
				Title:    "Transfer",
				Subtitle: "Transfer to real human",
				Buttons: []RichMessageButton{
					{
						Type:       "message",
						Text:       "transfer",
						PostbackId: "transfer_to_human",
						Value:      "transfer",
						UserIDS:    []string{},
					},
				},
			},
		},
	}
	return SendEventDto{
		ChatId: chatId,
		Event:  richMessage,
	}
}
