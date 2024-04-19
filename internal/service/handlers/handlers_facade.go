package handlers

import (
	"net/http"
	"responder/internal/service/responder"
)

type ResponderHandlersFacade struct {
	incomingEventHandler *incomingEventResponderHandler
}

func NewResponderHandlersFacade(responder responder.Responder, webhookSecrets []string) *ResponderHandlersFacade {
	return &ResponderHandlersFacade{
		incomingEventHandler: newIncomingEventResponderHandler(responder, webhookSecrets),
	}
}

func (f *ResponderHandlersFacade) HandleIncomingEvent() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(f.incomingEventHandler.Handle)
}
