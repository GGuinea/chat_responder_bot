package agent

import (
	"bytes"
	"encoding/json"
	"net/http"
	"responder/config"
	"responder/pkg/lc_api/common"
)

type LcAgentApi interface {
	SendEvent(event interface{}) error
	SetBotRoutingStatus(botId string, status string) error
	ListAgentsIdsForTransfer(chatId string) ([]string, error)
	TransferChat(chatId string, newAgentId string) error
}

type BasicAgentApi struct {
	cfg    *config.Config
	client *http.Client
}

func NewBasicAgentApi(cfg *config.Config) *BasicAgentApi {
	return &BasicAgentApi{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (ba *BasicAgentApi) SendEvent(eventData interface{}) error {
	url := buildSendEventURL(*ba.cfg.ChatAPIConfig)

	_, err := ba.send(eventData, url)

	if err != nil {
		return err
	}

	return nil
}

func (ba *BasicAgentApi) SetBotRoutingStatus(botId, status string) error {
	url := buildSetBotStatusURL(*ba.cfg.ChatAPIConfig)
	reqStruct := newSetRouteStatusRequest(status, botId)

	_, err := ba.send(reqStruct, url)
	if err != nil {
		return err
	}
	return nil
}

func (ba *BasicAgentApi) ListAgentsIdsForTransfer(chatId string) ([]string, error) {
	url := buildListAgentsForTransferURL(*ba.cfg.ChatAPIConfig)
	reqStruct := newListAgentsForTransferRequest(chatId)

	response, err := ba.send(reqStruct, url)
	if err != nil {
		return nil, err
	}

	var respDto []listAgentsForTransferResponse

	err = json.NewDecoder(response.Body).Decode(&respDto)
	if err != nil {
		return []string{}, err
	}

	return getStringArrayResponse(respDto), nil
}

func getStringArrayResponse(agents []listAgentsForTransferResponse) (resp []string) {
	for _, agent := range agents {
		resp = append(resp, agent.AgentId)
	}
	return
}

func (ba *BasicAgentApi) TransferChat(chatId, newAgentId string) error {
	url := buildTransferChatURL(*ba.cfg.ChatAPIConfig)
	reqStruct := newTransferToAgentRequest(chatId, newAgentId)

	_, err := ba.send(reqStruct, url)
	if err != nil {
		return err
	}

	return nil
}

func (ba *BasicAgentApi) send(requestData interface{}, url string) (*http.Response, error) {
	body, err := json.Marshal(requestData)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	if ba.cfg.UsePAT {
		request.Header.Set("Authorization", "Basic "+ba.cfg.PAT)
	} else {
		request.Header.Set("Authorization", "Bearer "+ba.cfg.OauthConfig.AccessToken)
	}

	if ba.cfg.UseBot {
		request.Header.Set("X-Author-Id", ba.cfg.BotId)
	}
	response, err := ba.client.Do(request)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, common.DecodeError(response.Body)
	}

	return response, nil
}

func buildTransferChatURL(cfg config.ChatAPI) string {
	return cfg.BaseURL + cfg.APIVersion + "/agent/action/transfer_chat"
}

func buildListAgentsForTransferURL(cfg config.ChatAPI) string {
	return cfg.BaseURL + cfg.APIVersion + "/agent/action/list_agents_for_transfer"
}

func buildSendEventURL(cfg config.ChatAPI) string {
	return cfg.BaseURL + cfg.APIVersion + "/agent/action/send_event"
}

func buildSetBotStatusURL(cfg config.ChatAPI) string {
	return cfg.BaseURL + cfg.APIVersion + "/agent/action/set_routing_status"
}
