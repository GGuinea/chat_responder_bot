package configuration

import (
	"bytes"
	"encoding/json"
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

	request, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	response, err := bc.Send(request)

	if err != nil {
		slog.Error("Cannot make request", err)
	}

	if response.StatusCode != 200 {
		slog.Info("Status code different than 200; ", slog.Any("statusCode", response.StatusCode))
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
	request.Header.Set("Authorization", "Basic "+ba.cfg.PAT)
	return ba.client.Do(request)

}

func buildCreateBotURL(cfg config.ChatAPI) string {
	return cfg.BaseURL + cfg.APIVersion + "/configuration/action/create_bot"
}
