package main

import (
	"bytes"
	"fmt"
	"net/http"
	"responder/config"
)

func main() {
	config := config.BuildConfig()
	listChats(config)
}

func listChats(cfg *config.Config) {
	url := buildListChatURL(*cfg.ChatAPIConfig)

	request, err := http.NewRequest("POST", url, bytes.NewReader([]byte(`{}`)))
	if err != nil {
		panic(err)
	}
	fmt.Println(url)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Basic " + cfg.PAT)
	client := &http.Client{}
	response, err := client.Do(request) 
	if err != nil {
		panic(err)
	}
	fmt.Println(response.Body)
}

func buildListChatURL(cfg config.ChatAPI) string {
	return cfg.BaseURL + cfg.APIVersion + "/agent/action/list_chats"
}
