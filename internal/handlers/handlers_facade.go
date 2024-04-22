package handlers

import (
	"net/http"
	"responder/config"
	"responder/internal/service"
	"responder/pkg/lc_api/auth"
)

type ResponderHandlersFacade struct {
	webhookHandler          *webhookEventHandler
	authCodeCallbackHandler *authCodeCallbackHandler
}

type ResponderHandlerFacadeDeps struct {
	Responder service.Responder
	Config    *config.Config
	AuthApi   auth.LcAuthApi
}

func NewResponderHandlersFacade(deps *ResponderHandlerFacadeDeps) *ResponderHandlersFacade {
	return &ResponderHandlersFacade{
		webhookHandler:          newWebhookEventHandler(deps.Responder, deps.Config.WebhooksSecrets),
		authCodeCallbackHandler: newAuthCodeCallbackHandler(deps.Config, deps.AuthApi),
	}
}

func (f *ResponderHandlersFacade) HandleNewWebhookEvent() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(f.webhookHandler.Handle)
}

func (f *ResponderHandlersFacade) HandleAuthCode() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(f.authCodeCallbackHandler.Handle)
}
