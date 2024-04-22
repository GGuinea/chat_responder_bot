package handlers

import (
	"net/http"
	"responder/config"
	"responder/internal/service/responder"
	"responder/pkg/lc_api/auth"
)

type ResponderHandlersFacade struct {
	webhookHandler    *webhookEventHandler
	authCodeCallbackHandler *authCodeCallbackHandler
}

func NewResponderHandlersFacade(responder responder.Responder, webhookSecrets []string, config *config.Config, authApi auth.LcAuthApi) *ResponderHandlersFacade {
	return &ResponderHandlersFacade{
		webhookHandler:    newWebhookEventHandler(responder, webhookSecrets),
		authCodeCallbackHandler: newAuthCodeCallbackHandler(config, authApi),
	}
}

func (f *ResponderHandlersFacade) HandleNewWebhookEvent() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(f.webhookHandler.Handle)
}

func (f *ResponderHandlersFacade) HandleAuthCode() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(f.authCodeCallbackHandler.Handle)
}
