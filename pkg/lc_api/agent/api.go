package agent

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"responder/config"
)

type LcAgentApi interface {
	SendEvent(interface{}) error
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
		slog.Error("Cannot make request", err)
	}

	return nil
}

func (ba *BasicAgentApi) Send(request *http.Request) (*http.Response, error) {
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Basic "+ba.cfg.PAT)
	return ba.client.Do(request)

}

func buildSendEventURL(cfg config.ChatAPI) string {
	return cfg.BaseURL + cfg.APIVersion + "/agent/action/send_event"
}

func buildListChatURL(cfg config.ChatAPI) string {
	return cfg.BaseURL + cfg.APIVersion + "/agent/action/list_chats"
}
