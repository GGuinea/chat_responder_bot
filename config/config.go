package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	PAT                     string
	ClientID                string
	BotId                   string
	UseBot                  bool
	UsePAT                  bool
	WebhooksSecrets         []string
	GracefulShutdownTimeout int
	ChatAPIConfig           *ChatAPI
	OauthConfig             *OAuth
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
		GracefulShutdownTimeout: getIntEnv("GRACEFUL_SHUTDOWN_TIMEOUT"),
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

func getIntEnv(name string) int {
	env := os.Getenv(name)
	parsed, err := strconv.Atoi(env)
	if err != nil {
		return 0
	}
	return parsed
}
