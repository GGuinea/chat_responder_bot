package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type LcAuthApi interface {
	Exchange(exchangeReq interface{}) (*AccessTokenGrantResponse, error)
}

type BasicAuthApi struct {
	client *http.Client
}

func NewBasichAuthApi() *BasicAuthApi {
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
		slog.Error(err.Error())
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
		slog.Error("Cannot make request", err)
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		errorBody, err := io.ReadAll(response.Body)
		if err != nil {
			slog.Error("Cannot decode body for error")
		}
		slog.Error("Status code different than 200; ", slog.Any("statusCode", response.StatusCode), slog.Any("error", string(errorBody)))
		return nil, fmt.Errorf("Status code different than 200; %v", response.StatusCode)
	}

	return response, nil
}

func getAccessTokenExchangeURL() string {
	return "https://accounts.livechat.com/v2/token"
}
