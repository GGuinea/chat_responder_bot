package main

import (
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

	agentApi := agent.NewBasicAgentApi(config)
	configurationApi := configuration.NewBasicConfiguratioApi(config)
	configurationApi.CreateBot(bots.NewDefaultBot("testowy", config.ClientID))

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
