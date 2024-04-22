package agent

type transferToAgentRequest struct {
	ChatId string         `json:"id"`
	Target transferTarget `json:"target"`
}

type transferTarget struct {
	Type string   `json:"type"`
	IDs  []string `json:"ids"`
}

func newTransferToAgentRequest(chatId, agentId string) transferToAgentRequest {
	transferTargetStruct := transferTarget{
		Type: "agent",
		IDs: []string{
			agentId,
		},
	}

	return transferToAgentRequest{
		ChatId: chatId,
		Target: transferTargetStruct,
	}
}
