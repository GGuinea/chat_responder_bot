package agent

type listAgentsForTransferRequest struct {
	ChatId string `json:"chat_id"`
}

type listAgentsForTransferResponse struct {
	AgentId string `json:"agent_id"`
}

func newListAgentsForTransferRequest(chatId string) listAgentsForTransferRequest {
	return listAgentsForTransferRequest{
		ChatId: chatId,
	}
}
