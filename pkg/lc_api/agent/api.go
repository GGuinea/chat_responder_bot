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

	request, err := http.NewRequest("POST", url, bytes.NewReader(body))
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

	request, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	_, err = ba.Send(request)
	if err != nil {
		return nil
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

	if response.StatusCode != 200 {
		errorBody, err := io.ReadAll(response.Body)
		if err != nil {
			slog.Error("Cannot decode body for error")
		}
		slog.Error("Status code different than 200; ", slog.Any("statusCode", response.StatusCode), slog.Any("error", string(errorBody)))
		return nil, fmt.Errorf("Status code different than 200; %v", response.StatusCode)
	}

	return response, nil

}

func buildSendEventURL(cfg config.ChatAPI) string {
	return cfg.BaseURL + cfg.APIVersion + "/agent/action/send_event"
}

func buildSetBotStatusURL(cfg config.ChatAPI) string {
	return cfg.BaseURL + cfg.APIVersion + "/agent/action/set_routing_status"
}
