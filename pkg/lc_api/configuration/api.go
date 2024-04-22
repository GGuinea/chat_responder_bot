package configuration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"responder/config"
	"responder/pkg/lc_api/common"
)

type LcConfigurationApi interface {
	CreateBot(botData interface{}) (*string, error)
	ListBots() ([]string, error)
}

type BasicConfigurationApi struct {
	cfg    *config.Config
	client *http.Client
}

func NewBasicConfiguratioApi(cfg *config.Config) *BasicConfigurationApi {
	return &BasicConfigurationApi{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (bc *BasicConfigurationApi) CreateBot(createBotData interface{}) (*string, error) {
	url := buildCreateBotURL(*bc.cfg.ChatAPIConfig)

	response, err := bc.send(createBotData, url)
	if err != nil {
		return nil, fmt.Errorf("Cannot create bot; %w", err)
	}

	var respDto CreateBotResponse

	err = json.NewDecoder(response.Body).Decode(&respDto)
	if err != nil {
		slog.Error("Cannot decode create bot response; ", err)
	}

	return &respDto.Id, nil
}

func (bc *BasicConfigurationApi) ListBots() ([]string, error) {
	url := buildListBotsURL(*bc.cfg.ChatAPIConfig)
	response, err := bc.send(newListBotsRequest(), url)

	if err != nil {
		return nil, fmt.Errorf("Cannot create bot; %w", err)
	}

	var respDto []listBotsResponse

	err = json.NewDecoder(response.Body).Decode(&respDto)

	if err != nil {
		slog.Error("Cannot decode create bot response; ", err)
	}
	return extractIds(respDto), nil
}

func extractIds(resp []listBotsResponse) (ids []string) {
	for _, bot := range resp {
		ids = append(ids, bot.Id)
	}
	return
}

func (ba *BasicConfigurationApi) send(requestData interface{}, url string) (*http.Response, error) {
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

	response, err := ba.client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, common.DecodeError(response.Body)
	}

	return response, nil
}

func buildCreateBotURL(cfg config.ChatAPI) string {
	return cfg.BaseURL + cfg.APIVersion + "/configuration/action/create_bot"
}

func buildListBotsURL(cfg config.ChatAPI) string {
	return cfg.BaseURL + cfg.APIVersion + "/configuration/action/list_bots"
}
