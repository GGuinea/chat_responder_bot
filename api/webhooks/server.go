package webhooks

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"responder/config"
)

type HttpServer interface {
	Start() error
	GracefulStop(ctx context.Context) error
}

type WebhookServer struct {
	server *http.Server
}

type WebhookServerDeps struct {
	config *config.Config
}

func NewWebhookServer(deps *WebhookServerDeps) (HttpServer, error) {
	if deps == nil {
		return nil, errors.New("Deps cannot be null")
	}

	mux := http.NewServeMux()
	buildRouter(mux)

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	return &WebhookServer{server: server}, nil
}

func buildRouter(mux *http.ServeMux) {
	mux.HandleFunc("GET /ping", pong)
	mux.HandleFunc("POST /ping", pong)
}

func pong(w http.ResponseWriter, r *http.Request) {
	slog.Info("received")
	w.Write([]byte("pong"))
}

func (whs *WebhookServer) Start() error {
	return whs.server.ListenAndServe()
}

func (whs *WebhookServer) GracefulStop(ctx context.Context) error {
	return whs.server.Shutdown(ctx)
}
