package service

import (
	"fmt"
	"log/slog"
	"responder/config"
	"responder/internal/model"
	"responder/pkg/lc_api/agent"
	"responder/pkg/lc_api/configuration"
	"slices"
)

type Responder interface {
	Start()
	GracefulStop()
	HandleIncomingEvent(event model.IncomingEvent)
	HandleRichMessagePostback(event model.RichMessagePostbackEvent)
}

type BasicResponder struct {
	incomingEvents   chan model.ResponderEvent
	chatApi          agent.LcAgentApi
	configurationAPI configuration.LcConfigurationApi
	close            chan struct{}
	config           *config.Config
}

type ResponderDeps struct {
	ChatApi          agent.LcAgentApi
	Config           *config.Config
	ConfigurationApi configuration.LcConfigurationApi
}

func NewResponder(deps *ResponderDeps) *BasicResponder {
	return &BasicResponder{
		incomingEvents:   make(chan model.ResponderEvent, 20),
		chatApi:          deps.ChatApi,
		configurationAPI: deps.ConfigurationApi,
		close:            make(chan struct{}),
		config:           deps.Config,
	}
}

func (r *BasicResponder) Start() {
	for {
		select {
		case event, ok := <-r.incomingEvents:
			if !ok {
				slog.Error("Problem while reading from input channel;")
				break
			}

			if err := r.doResponse(event); err != nil {
				slog.Error("Cannot make response; ", err)
			}
		case <-r.close:
			return
		}
	}
}

func (r *BasicResponder) HandleIncomingEvent(event model.IncomingEvent) {
	//r.incomingEvents <- model.NewPlainMessageResponderEvent(event.Payload.ChatId)
	r.incomingEvents <- model.NewRichMessageResponderEvent(event.Payload.ChatId)
}

func (r *BasicResponder) HandleRichMessagePostback(event model.RichMessagePostbackEvent) {
	if event.Payload.Postback.ActionId == model.ACTION_NO {
		r.incomingEvents <- model.NewPlainMessageResponderEvent(event.Payload.ChatId)
	} else {
		r.incomingEvents <- model.NewTransferResponderEvent(event.Payload.ChatId)
	}
}

func (r *BasicResponder) doResponse(event model.ResponderEvent) error {
	slog.Info("Trying to send response")
	var response model.SendEventDto
	switch event.ActionToPerform {
	case model.PLAIN_MESSAGE_REPLY:
		response = model.NewDefaultMessageEvent(event.ChatId, "plain text response")
	case model.RICH_MESSAGE_REPLY:
		response = model.NewRichCardMessageEvent(event.ChatId)
	case model.TRANSFER_TO_HUMAN_AGENT:
		err := r.transferToHuman(event.ChatId)
		if err != nil {
			slog.Warn("Cannot transfer", slog.Any("msg", err))
			response = model.NewDefaultMessageEvent(event.ChatId, "cannot transfer right now, please try again")
		} else {
			return nil
		}
	default:
		return fmt.Errorf("Unknow reply type; got: %v", event.ActionToPerform)
	}
	return r.chatApi.SendEvent(response)
}

func (r *BasicResponder) GracefulStop() {
	r.close <- struct{}{}
}

func (r *BasicResponder) transferToHuman(chatId string) error {
	availableAgents, err := r.chatApi.ListAgentsIdsForTransfer(chatId)
	if err != nil {
		return err
	}

	botsIds, err := r.configurationAPI.ListBots()
	if err != nil {
		return fmt.Errorf("Cannot list bots; %w", err)
	}

	filtered := filterOutBots(availableAgents, botsIds)

	slog.Info("Transfer summary", slog.Any("availableAgents", len(availableAgents)), slog.Any("bots", len(botsIds)), slog.Any("after filtration", len(filtered)))
	if len(filtered) == 0 {
		return fmt.Errorf("No available human agents!")
	}

	err = r.chatApi.TransferChat(chatId, filtered[0])
	if err != nil {
		return fmt.Errorf("Cannot transfer %v", err)
	}

	return nil
}

func filterOutBots(availableAgentsIds []string, botsIds []string) (properIds []string) {
	for _, agentId := range availableAgentsIds {
		if !slices.Contains(botsIds, agentId) {
			properIds = append(properIds, agentId)
		}
	}
	return
}
