package main

import (
	"responder/api/webhooks"
	"responder/config"
	"responder/internal/model"
	"responder/internal/service/handlers"
	"responder/internal/service/responder"
	lcapi "responder/pkg/lc_api"
)

func main() {
	config := config.BuildConfig()

	api := lcapi.NewBasicApi(config)
	incomingEventsCh := make(chan model.IncomingEvent, 20)
	responderDeps := responder.ResponderDeps{
		IncomingEventsCh: incomingEventsCh,
		ChatApi:          api,
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
	for {}
}
