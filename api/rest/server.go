package rest

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

type RestServer struct {
	server *http.Server
}

type RestServerDeps struct {
	HandlersFacade *handlers.ResponderHandlersFacade
}

func NewRestServer(deps *RestServerDeps) (HttpServer, error) {
	if deps == nil {
		return nil, errors.New("Deps cannot be null")
	}

	mux := http.NewServeMux()
	buildRoutes(mux, deps.HandlersFacade)

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: mux,
	}

	return &RestServer{server: server}, nil
}

func buildRoutes(mux *http.ServeMux, handlerFacade *handlers.ResponderHandlersFacade) {
	mux.HandleFunc("GET /ping", pong)
	mux.HandleFunc("POST /webhook", handlerFacade.HandleNewWebhookEvent())
	mux.HandleFunc("GET /auth", handlerFacade.HandleAuthCode())
}

func pong(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

func (whs *RestServer) Start() error {
	return whs.server.ListenAndServe()
}

func (whs *RestServer) GracefulStop(ctx context.Context) error {
	return whs.server.Shutdown(ctx)
}
