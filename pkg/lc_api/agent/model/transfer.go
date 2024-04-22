package model

type TransferToAgentRequest struct {
	ChatId string         `json:"id"`
	Target transferTarget `json:"target"`
}

type transferTarget struct {
	Type string   `json:"type"`
	IDs  []string `json:"ids"`
}

func NewTransferToAgentRequest(chatId, agentId string) TransferToAgentRequest {
	transferTargetStruct := transferTarget{
		Type: "agent",
		IDs: []string{
			agentId,
		},
	}

	return TransferToAgentRequest{
		ChatId: chatId,
		Target: transferTargetStruct,
	}
}
