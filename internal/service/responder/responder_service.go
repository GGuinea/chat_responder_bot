package responder

import (
	"fmt"
	"log/slog"
	"responder/internal/model"
	"responder/pkg/lc_api/agent"
)

type Responder interface {
	Start()
	GracefulStop()
	HandleIncomingEvent(event model.IncomingEvent)
	HandleRichMessagePostback(event model.RichMessagePostbackEvent)
}

type BasicResponder struct {
	incomingEvents chan model.ResponderEvent
	chatApi        agent.LcAgentApi
	close          chan struct{}
}

type ResponderDeps struct {
	IncomingEventsCh chan model.ResponderEvent
	ChatApi          agent.LcAgentApi
}

func NewResponder(deps *ResponderDeps) *BasicResponder {
	return &BasicResponder{
		incomingEvents: make(chan model.ResponderEvent, 20),
		chatApi:        deps.ChatApi,
	}
}

func (r *BasicResponder) Start() {
	for {
		select {
		case event, ok := <-r.incomingEvents:
			if !ok {
				slog.Error("Problem while reading from input channel;")
				return
			}

			if err := r.doResponse(event); err != nil {
				slog.Error("Cannot make response; ", err)
				return
			}
		case <-r.close:
			return
		}
	}
}

func (r *BasicResponder) HandleIncomingEvent(event model.IncomingEvent) {
	r.incomingEvents <- model.NewPlainMessageResponderEvent(event.Payload.ChatId)
}

func (r *BasicResponder) HandleRichMessagePostback(event model.RichMessagePostbackEvent) {
	r.incomingEvents <- model.NewRichMessageResponderEvent(event.Payload.ChatId)
}

func (r *BasicResponder) doResponse(event model.ResponderEvent) error {
	slog.Info("Trying to send response")
	var response model.SendEventDto
	switch event.ReplyType {
	case model.PLAIN_MESSAGE_REPLY:
		response = model.NewDefaultMessageEvent(event.ChatId, "plain text response")
	case model.RICH_MESSAGE_REPLY:
		response = model.NewRichCardMessageEvent(event.ChatId)
	default:
		return fmt.Errorf("Unknow reply type; got: %v", event.ReplyType)
	}
	return r.chatApi.SendEvent(response)
}

func (r *BasicResponder) GracefulStop() {
	r.close <- struct{}{}
}
