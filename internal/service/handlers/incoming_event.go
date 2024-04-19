package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"responder/internal/model"
	"responder/internal/service/responder"
	"slices"
)

type incomingEventResponderHandler struct {
	responder responder.Responder
	secrets   []string
}

func newIncomingEventResponderHandler(responder responder.Responder, secrets []string) *incomingEventResponderHandler {
	return &incomingEventResponderHandler{
		responder: responder,
		secrets:   secrets,
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

	fmt.Printf("%#+v\n", incomingEvent)
	if !slices.Contains(ierh.secrets, incomingEvent.SecretKey) {
		slog.Warn("Cannot process incoming event, wrong secret")
	}

	ierh.responder.HandleNewEvent(incomingEvent)
}
