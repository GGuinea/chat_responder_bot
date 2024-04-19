package handlers

import (
	"net/http"
	"responder/config"
	"responder/internal/service/responder"
	"responder/pkg/lc_api/auth"
)

type ResponderHandlersFacade struct {
	incomingEventHandler    *incomingEventResponderHandler
	authCodeCallbackHandler *authCodeCallbackHandler
}

func NewResponderHandlersFacade(responder responder.Responder, webhookSecrets []string, config *config.Config, authApi auth.LcAuthApi) *ResponderHandlersFacade {
	return &ResponderHandlersFacade{
		incomingEventHandler:    newIncomingEventResponderHandler(responder, webhookSecrets),
		authCodeCallbackHandler: newAuthCodeCallbackHandler(config, authApi),
	}
}

func (f *ResponderHandlersFacade) HandleIncomingEvent() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(f.incomingEventHandler.Handle)
}

func (f *ResponderHandlersFacade) HandleAuthCode() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(f.authCodeCallbackHandler.Handle)
}
