package main

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rifat-simoom/go-clean-architecture/internal/common/logs"
	"github.com/rifat-simoom/go-clean-architecture/internal/common/server"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainings/ports"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainings/service"
)

func main() {
	logs.Init()

	ctx := context.Background()

	app, cleanup := service.NewApplication(ctx)
	defer cleanup()

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return ports.HandlerFromMux(ports.NewHttpServer(app), router)
	})
}
