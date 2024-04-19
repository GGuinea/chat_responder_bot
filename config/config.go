package config

import (
	"os"
	"strings"
	"time"
)

type Config struct {
	PAT             string
	ClientID        string
	BotId           string
	UseBot          bool
	UsePAT          bool
	WebhooksSecrets []string
	ChatAPIConfig   *ChatAPI
	OauthConfig     *OAuth
}

type ChatAPI struct {
	BaseURL    string
	APIVersion string
}

type OAuth struct {
	ClientSecret        string
	AccessToken         string
	RefreshToken        string
	ExpiresIn           int
	RedirectURI         string
	AccsesTokenCreation time.Time
}

func BuildConfig() *Config {
	return &Config{
		PAT:             os.Getenv("PAT"),
		ClientID:        os.Getenv("CLIENT_ID"),
		WebhooksSecrets: parseWebhooksecrets(),
		ChatAPIConfig:   buildChatApi(),
		OauthConfig:     buildOAuth(),
	}
}

func buildChatApi() *ChatAPI {
	return &ChatAPI{
		BaseURL:    os.Getenv("API_BASE_URL"),
		APIVersion: os.Getenv("API_VERSION"),
	}
}

func buildOAuth() *OAuth {
	return &OAuth{
		ClientSecret:        os.Getenv("CLIENT_SECRET"),
		AccessToken:         os.Getenv("ACCESS_TOKEN"),
		RefreshToken:        "",
		ExpiresIn:           0,
		RedirectURI:         os.Getenv("REDIRECT_URI"),
		AccsesTokenCreation: time.Time{},
	}
}

func (c *Config) SetUseBotFlag(val bool) {
	c.UseBot = val
}

func (c *Config) SetBotId(val string) {
	c.BotId = val
}

func (c *Config) SetUsePATFlag(val bool) {
	c.UsePAT = val
}

func (c *Config) SetAccessToken(val string) {
	c.OauthConfig.AccessToken = val
}

func (c *Config) SetRefreshToken(val string) {
	c.OauthConfig.RefreshToken = val
}

func (c *Config) SetExpiresIn(val int) {
	c.OauthConfig.ExpiresIn = val
}

func (c *Config) SetAccessTokenCreationTime(val time.Time) {
	c.OauthConfig.AccsesTokenCreation = val
}

func parseWebhooksecrets() []string {
	stringWithCommas := os.Getenv("WEBHOOK_SECRETS")
	return strings.Split(stringWithCommas, ",")
}
