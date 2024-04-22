package agent

type setRouteStatusRequest struct {
	Status  string `json:"status"`
	AgentId string `json:"agent_id"`
}

func newSetRouteStatusRequest(status, agentId string) setRouteStatusRequest {
	return setRouteStatusRequest{
		Status:  status,
		AgentId: agentId,
	}
}
