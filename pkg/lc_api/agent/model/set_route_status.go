package model

type SetRouteStatusRequest struct {
	Status  string `json:"status"`
	AgentId string `json:"agent_id"`
}

func NewSetRouteStatusRequest(status, agentId string) SetRouteStatusRequest {
	return SetRouteStatusRequest{
		Status:  status,
		AgentId: agentId,
	}
}
