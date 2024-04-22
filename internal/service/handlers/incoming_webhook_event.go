package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"responder/internal/model"
	"responder/internal/service/responder"
	"slices"
)

type webhookEventHandler struct {
	responder responder.Responder
	secrets   []string
}

func newWebhookEventHandler(responder responder.Responder, secrets []string) *webhookEventHandler {
	return &webhookEventHandler{
		responder: responder,
		secrets:   secrets,
	}
}

func (weh *webhookEventHandler) Handle(_ http.ResponseWriter, r *http.Request) {
	slog.Info("New event received")
	var webhookEvent model.WebhookEvent

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Cannot read body bytes;", err)
		return
	}

	err = json.Unmarshal(body, &webhookEvent)
	if err != nil {
		slog.Error("Cannot unmarshal body;", err)
		return
	}

	if !slices.Contains(weh.secrets, webhookEvent.SecretKey) {
		slog.Warn("Cannot process incoming event, wrong secret key")
		return
	}
	weh.handleEvent(webhookEvent.Action, body)
}

func (weh *webhookEventHandler) handleEvent(action string, eventData []byte) {
	switch action {
	case "incoming_event":
		var incomingEvent model.IncomingEvent
		err := json.Unmarshal(eventData, &incomingEvent)
		if err != nil {
			slog.Error("Cannot unmarshal body")
			return
		}

		fmt.Printf("Incoming event received: %#+v\n", incomingEvent)
		weh.responder.HandleIncomingEvent(incomingEvent)
	case "incoming_rich_message_postback":
		var richMessagePostback model.RichMessagePostbackEvent
		err := json.Unmarshal(eventData, &richMessagePostback)
		if err != nil {
			slog.Error("Cannot unmarshal body")
			return
		}

		fmt.Printf("Rich message postback received: %#+v\n", richMessagePostback)
		weh.responder.HandleRichMessagePostback(richMessagePostback)
	default:
		slog.Warn("Unknow action", slog.Any("action", action))
		return
	}
}
