package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"responder/internal/model"
	"responder/internal/service/responder"
)

type incomingEventResponderHandler struct {
	responder responder.Responder
}

func newIncomingEventResponderHandler(responder responder.Responder) *incomingEventResponderHandler {
	return &incomingEventResponderHandler{
		responder: responder,
	}
}

func (ierh *incomingEventResponderHandler) Handle(_ http.ResponseWriter, r *http.Request) {
	slog.Info("event received")
	var incomingEvent model.IncomingEvent

	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&incomingEvent)
	if err != nil {
		slog.Info("Cannot decode body;", err)
		return
	}

	ierh.responder.HandleNewEvent(incomingEvent)
}
