package main

import (
	"flag"
	"fmt"
	"log/slog"
	"math/rand"
	"responder/api/webhooks"
	"responder/config"
	"responder/internal/model"
	"responder/internal/model/bots"
	"responder/internal/service/handlers"
	"responder/internal/service/responder"
	"responder/pkg/lc_api/agent"
	"responder/pkg/lc_api/configuration"
)

func main() {
	config := config.BuildConfig()
	var useBot bool
	var botId string

	flag.BoolVar(&useBot, "use_bot", false, "Ans as bot")
	flag.StringVar(&botId, "bot_id", "", "Bot id to ans")
	flag.Parse()

	if useBot && botId == "" {
		slog.Info("Creating new bot")
		botId = createNewBot(config)
	}

	config.SetUseBotFlag(useBot)
	config.SetBotId(botId)

	agentApi := agent.NewBasicAgentApi(config)

	incomingEventsCh := make(chan model.IncomingEvent, 20)
	responderDeps := responder.ResponderDeps{
		IncomingEventsCh: incomingEventsCh,
		ChatApi:          agentApi,
	}
	responder := responder.NewResponder(&responderDeps)
	handlersFacade := handlers.NewResponderHandlersFacade(responder)

	webhookServer, err := webhooks.NewWebhookServer(&webhooks.WebhookServerDeps{
		HandlersFacade: handlersFacade,
	})

	if err != nil {
		panic(err)
	}

	go webhookServer.Start()
	go responder.Start()
	for {
	}
}

func createNewBot(config *config.Config) string {
	configurationApi := configuration.NewBasicConfiguratioApi(config)
	botId, err := configurationApi.CreateBot(bots.NewDefaultBot(fmt.Sprintf("testowy %d", rand.Int()), config.ClientID))
	if err != nil {
		slog.Error("Cannot create new bot", err)
		panic(err)
	}

	return *botId
}
