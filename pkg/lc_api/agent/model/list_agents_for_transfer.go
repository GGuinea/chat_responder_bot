package model

type ListAgentsForTransferRequest struct {
	ChatId string `json:"chat_id"`
}

type ListAgentsForTransferResponse struct {
	AgentId string `json:"agent_id"`
}

func NewListAgentsForTransferRequest(chatId string) ListAgentsForTransferRequest {
	return ListAgentsForTransferRequest{
		ChatId: chatId,
	}
}
