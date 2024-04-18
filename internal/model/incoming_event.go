package model

type IncomingEventPayload struct {
	ChatId   string `json:"chat_id"`
	ThreadId string `json:"thread_id"`
}

type IncomingEvent struct {
	WebhookId      string               `json:"webhook_id"`
	SecretKey      string               `json:"secret_key"` // how to verify it?
	Action         string               `json:"action"`
	OrganizationId string               `json:"organization_id"`
	Payload        IncomingEventPayload `json:"payload"`
}
