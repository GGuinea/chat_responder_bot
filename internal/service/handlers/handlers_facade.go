package handlers

import (
	"net/http"
	"responder/internal/service/responder"
)

type ResponderHandlersFacade struct {
	incomingEventHandler *incomingEventResponderHandler
}

func NewResponderHandlersFacade(responder responder.Responder) *ResponderHandlersFacade {
	return &ResponderHandlersFacade{
		incomingEventHandler: newIncomingEventResponderHandler(responder),
	}
}

func (f *ResponderHandlersFacade) HandleIncomingEvent() func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(f.incomingEventHandler.Handle)
}
