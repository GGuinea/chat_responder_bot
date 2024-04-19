package responder

import (
	"log/slog"
	"responder/internal/model"
	"responder/pkg/lc_api/agent"
)

type Responder interface {
	Start()
	HandleNewEvent(event model.IncomingEvent)
	GracefulStop()
}

type BasicResponder struct {
	incomingEvents chan model.IncomingEvent
	chatApi        agent.LcAgentApi
	close          chan struct{}
}

type ResponderDeps struct {
	IncomingEventsCh chan model.IncomingEvent
	ChatApi          agent.LcAgentApi
}

func NewResponder(deps *ResponderDeps) *BasicResponder {
	return &BasicResponder{
		incomingEvents: deps.IncomingEventsCh,
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

func (r *BasicResponder) HandleNewEvent(event model.IncomingEvent) {
	r.incomingEvents <- event
}

func (r *BasicResponder) doResponse(event model.IncomingEvent) error {
	slog.Info("Trying to send event")
	response := model.NewRichCardMessageEvent(event.Payload.ChatId)
	return r.chatApi.SendEvent(response)
}

func (r *BasicResponder) GracefulStop() {
	r.close <- struct{}{}
}
