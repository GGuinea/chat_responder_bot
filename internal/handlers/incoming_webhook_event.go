package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"responder/internal/model"
	"responder/internal/service"
	"slices"
)

type webhookEventHandler struct {
	responder service.Responder
	secrets   []string
}

func newWebhookEventHandler(responder service.Responder, secrets []string) *webhookEventHandler {
	return &webhookEventHandler{
		responder: responder,
		secrets:   secrets,
	}
}

var (
	EVENT_TYPE_INCOMING_EVENT                 = "incoming_event"
	EVENT_TYPE_INCOMING_RICH_MESSAGE_POSTBACK = "incoming_rich_message_postback"
)

func (weh *webhookEventHandler) Handle(_ http.ResponseWriter, r *http.Request) {
	slog.Info("New event received")
	var webhookEvent model.WebhookMsg

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
	case EVENT_TYPE_INCOMING_EVENT:
		var incomingEvent model.IncomingEvent
		err := json.Unmarshal(eventData, &incomingEvent)
		if err != nil {
			slog.Error("Cannot unmarshal body")
			return
		}

		fmt.Printf("Incoming event received: %#+v\n", incomingEvent)
		weh.responder.HandleIncomingEvent(incomingEvent)
	case EVENT_TYPE_INCOMING_RICH_MESSAGE_POSTBACK:
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
