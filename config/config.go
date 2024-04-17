package config

import "os"

type Config struct {
	PAT           string
	ChatAPIConfig *ChatAPI
}

type ChatAPI struct {
	BaseURL    string
	APIVersion string
}

func BuildConfig() *Config {
	return &Config{
		PAT:           os.Getenv("PAT"),
		ChatAPIConfig: buildChatApi(),
	}
}

func buildChatApi() *ChatAPI {
	return &ChatAPI{
		BaseURL:    os.Getenv("API_BASE_URL"),
		APIVersion: os.Getenv("API_VERSION"),
	}
}
