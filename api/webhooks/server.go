package webhooks

import (
	"context"
	"errors"
	"net/http"
	"responder/internal/service/handlers"
)

type HttpServer interface {
	Start() error
	GracefulStop(ctx context.Context) error
}

type WebhookServer struct {
	server *http.Server
}

type WebhookServerDeps struct {
	HandlersFacade *handlers.ResponderHandlersFacade
}

func NewWebhookServer(deps *WebhookServerDeps) (HttpServer, error) {
	if deps == nil {
		return nil, errors.New("Deps cannot be null")
	}

	mux := http.NewServeMux()
	buildRoutes(mux, deps.HandlersFacade)

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	return &WebhookServer{server: server}, nil
}

func buildRoutes(mux *http.ServeMux, handlerFacade *handlers.ResponderHandlersFacade){
	mux.HandleFunc("GET /ping", pong)
	mux.HandleFunc("POST /incoming_event", handlerFacade.HandleIncomingEvent())
}

func pong(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func (whs *WebhookServer) Start() error {
	return whs.server.ListenAndServe()
}

func (whs *WebhookServer) GracefulStop(ctx context.Context) error {
	return whs.server.Shutdown(ctx)
}
