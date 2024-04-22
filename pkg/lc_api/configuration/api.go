package configuration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"responder/config"
	"responder/pkg/lc_api/configuration/model"
)

type LcConfigurationApi interface {
	CreateBot() (string, error)
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

	body, err := json.Marshal(createBotData)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	response, err := bc.Send(request)
	if err != nil {
		return nil, fmt.Errorf("Cannot create bot")
	}

	var respDto model.CreateBotResponse

	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&respDto)
	if err != nil {
		slog.Error("Cannot decode create bot response; ", err)
	}

	slog.Info(respDto.Id)
	return &respDto.Id, nil
}

func (ba *BasicConfigurationApi) Send(request *http.Request) (*http.Response, error) {
	request.Header.Set("Content-Type", "application/json")
	if ba.cfg.UsePAT {
		request.Header.Set("Authorization", "Basic "+ba.cfg.PAT)
	} else {
		request.Header.Set("Authorization", "Bearer "+ba.cfg.OauthConfig.AccessToken)
	}

	response, err := ba.client.Do(request)
	if err != nil {
		slog.Error("Cannot make request", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		errorBody, err := io.ReadAll(response.Body)
		if err != nil {
			slog.Error("Cannot decode body for error")
		}
		slog.Error("Status code different than 200; ", slog.Any("statusCode", response.StatusCode), slog.Any("error", string(errorBody)))
	}

	return response, nil
}

func buildCreateBotURL(cfg config.ChatAPI) string {
	return cfg.BaseURL + cfg.APIVersion + "/configuration/action/create_bot"
}
