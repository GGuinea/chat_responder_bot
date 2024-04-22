package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"responder/config"
	"responder/pkg/lc_api/agent/model"
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
	body, err := json.Marshal(eventData)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	_, err = ba.Send(request)

	if err != nil {
		return err
	}

	return nil
}

func (ba *BasicAgentApi) SetBotRoutingStatus(botId, status string) error {
	url := buildSetBotStatusURL(*ba.cfg.ChatAPIConfig)
	reqStruct := model.NewSetRouteStatusRequest(status, botId)

	body, err := json.Marshal(reqStruct)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	_, err = ba.Send(request)
	if err != nil {
		return err
	}
	return nil
}

func (ba *BasicAgentApi) ListAgentsIdsForTransfer(chatId string) ([]string, error) {
	url := buildListAgentsForTransferURL(*ba.cfg.ChatAPIConfig)
	reqStruct := model.NewListAgentsForTransferRequest(chatId)

	body, err := json.Marshal(reqStruct)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	response, err := ba.Send(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, decodeErrorMsg(response.Body)
	}

	var respDto []model.ListAgentsForTransferResponse

	err = json.NewDecoder(response.Body).Decode(&respDto)
	if err != nil {
		slog.Error("Cannot decode create bot response; ", err)
	}

	return getStringArrayResponse(respDto), nil
}

func getStringArrayResponse(agents []model.ListAgentsForTransferResponse) (resp []string) {
	for _, agent := range agents {
		resp = append(resp, agent.AgentId)
	}
	return
}

func (ba *BasicAgentApi) TransferChat(chatId, newAgentId string) error {
	url := buildTransferChatURL(*ba.cfg.ChatAPIConfig)
	reqStruct := model.NewTransferToAgentRequest(chatId, newAgentId)

	body, err := json.Marshal(reqStruct)
	if err != nil {
		return err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	response, err := ba.Send(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return decodeErrorMsg(response.Body)
	}

	return nil
}

func (ba *BasicAgentApi) Send(request *http.Request) (*http.Response, error) {
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
		slog.Error("Cannot make request", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, decodeErrorMsg(response.Body)
	}

	return response, nil

}

func decodeErrorMsg(body io.ReadCloser) error {
	errorBody, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("Cannot even decode body for error")
	}
	return fmt.Errorf("Status code different than 200, %s", string(errorBody))
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
