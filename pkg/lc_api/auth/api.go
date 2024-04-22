package auth

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"responder/pkg/lc_api/common"
)

type LcAuthApi interface {
	Exchange(exchangeReq interface{}) (*AccessTokenGrantResponse, error)
}

type BasicAuthApi struct {
	client *http.Client
}

func NewBasicAuthApi() *BasicAuthApi {
	return &BasicAuthApi{
		client: &http.Client{},
	}
}

func (baa *BasicAuthApi) Exchange(exchangeReq interface{}) (*AccessTokenGrantResponse, error) {
	baseUrl := getAccessTokenExchangeURL()

	body, err := json.Marshal(exchangeReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, baseUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	response, err := send(req)
	if err != nil {
		return nil, err
	}

	var accessTokenReponse AccessTokenGrantResponse

	err = json.NewDecoder(response.Body).Decode(&accessTokenReponse)
	if err != nil {
		return nil, err
	}

	return &accessTokenReponse, nil
}

func send(request *http.Request) (*http.Response, error) {
	slog.Info("sending")
	request.Header.Set("Content-Type", "application/json")
	client := http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, common.DecodeError(response.Body)
	}

	return response, nil
}

func getAccessTokenExchangeURL() string {
	return "https://accounts.livechat.com/v2/token"
}
