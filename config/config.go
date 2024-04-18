package config

import "os"

type Config struct {
	PAT           string
	ClientID      string
	BotId         string
	UseBot        bool
	ChatAPIConfig *ChatAPI
}

type ChatAPI struct {
	BaseURL    string
	APIVersion string
}

func BuildConfig() *Config {
	return &Config{
		PAT:           os.Getenv("PAT"),
		ClientID:      os.Getenv("CLIENT_ID"),
		ChatAPIConfig: buildChatApi(),
	}
}

func buildChatApi() *ChatAPI {
	return &ChatAPI{
		BaseURL:    os.Getenv("API_BASE_URL"),
		APIVersion: os.Getenv("API_VERSION"),
	}
}

func (c *Config) SetUseBotFlag(val bool) {
	c.UseBot = val
}

func (c *Config) SetBotId(val string) {
	c.BotId = val
}
