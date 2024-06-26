package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"responder/api/rest"
	"responder/config"
	"responder/internal/handlers"
	"responder/internal/service"
	"responder/pkg/lc_api/agent"
	"responder/pkg/lc_api/auth"
	"responder/pkg/lc_api/configuration"
	"syscall"
	"time"
)

func main() {
	config := config.BuildConfig()
	var useBot bool
	var usePAT bool
	var botId string

	flag.BoolVar(&useBot, "use_bot", false, "Ans as bot")
	flag.StringVar(&botId, "bot_id", "", "Bot id to ans")
	flag.BoolVar(&usePAT, "use_pat", false, "Use PAT with basic auth")

	flag.Parse()

	configurationApi := configuration.NewBasicConfiguratioApi(config)

	if useBot && botId == "" {
		slog.Info("Creating new bot...")
		botId = createNewBot(configurationApi, config)
	}

	config.SetUseBotFlag(useBot)
	config.SetBotId(botId)
	config.SetUsePATFlag(usePAT)

	agentApi := agent.NewBasicAgentApi(config)
	activateBot(agentApi, config.BotId)

	responder := service.NewResponder(&service.ResponderDeps{
		ChatApi:          agentApi,
		Config:           config,
		ConfigurationApi: configurationApi,
	})

	authApi := auth.NewBasicAuthApi()

	handlersFacade := handlers.NewResponderHandlersFacade(&handlers.ResponderHandlerFacadeDeps{
		Responder: responder,
		Config:    config,
		AuthApi:   authApi,
	})

	restServer, err := rest.NewRestServer(&rest.RestServerDeps{
		HandlersFacade: handlersFacade,
	})

	if err != nil {
		panic(err)
	}

	go restServer.Start()
	go responder.Start()

	stopCh := make(chan os.Signal)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-stopCh

	contextWithTimeout, cancel := context.WithTimeout(context.Background(), time.Duration(config.GracefulShutdownTimeout)*time.Second)
	defer cancel()

	if err := restServer.GracefulStop(contextWithTimeout); err != nil {
		slog.Info("Server shoutdown error", err)
	}

	select {
	case <-contextWithTimeout.Done():
		slog.Info("webhook server shoutdown, waiting for responder...")
		time.Sleep(time.Duration(config.GracefulShutdownTimeout) * time.Second)
		slog.Info("responder shoutdown, deactivating bot...")
		setNotAcceptChatsFlag(agentApi, config.BotId)
	}
}

func createNewBot(confAPI configuration.LcConfigurationApi, config *config.Config) string {
	botId, err := confAPI.CreateBot(configuration.NewDefaultBot(fmt.Sprintf("testowy %d", rand.Int()), config.ClientID))
	if err != nil {
		slog.Error("Cannot create new bot", err)
		panic(err)
	}

	return *botId
}

func activateBot(agentApi agent.LcAgentApi, botId string) error {
	return agentApi.SetBotRoutingStatus(botId, configuration.AGENT_STATUS_ACCEPTING_CHATS)
}

func setNotAcceptChatsFlag(agentApi agent.LcAgentApi, botId string) error {
	return agentApi.SetBotRoutingStatus(botId, configuration.AGENT_STATUS_NOT_ACCEPTING_CHATS)
}
