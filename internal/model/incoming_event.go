package model

type IncomingEventPayload struct {
	ChatId   string `json:"chat_id"`
	ThreadId string `json:"thread_id"`
	Event    Event  `json:"event"`
}

type IncomingEvent struct {
	WebhookId string               `json:"webhook_id"`
	SecretKey string               `json:"secret_key"`
	Action    string               `json:"action"`
	Payload   IncomingEventPayload `json:"payload"`
}

type Event struct {
	Type string `json:"type"`
}
