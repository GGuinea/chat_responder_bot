package model

type WebhookEvent struct {
	SecretKey string `json:"secret_key"`
	Action    string `json:"action"`
}

type IncomingEvent struct {
	Payload IncomingEventPayload `json:"payload"`
}

type IncomingEventPayload struct {
	ChatId   string `json:"chat_id"`
	ThreadId string `json:"thread_id"`
	Event    Event  `json:"event"`
}

type Event struct {
	Type string `json:"type"`
}

type RichMessagePostbackEvent struct {
	Payload IncomingEventPayload `json:"payload"`
}

type RichMessagePostbackEventPayload struct {
	ChatId string `json:"chat_id"`
}

type RichMessagePostback struct {
	ActionId string `json:"id"`
}
