package main

import (
	"context"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainings/infrastructure/configs"
	http2 "github.com/rifat-simoom/go-clean-architecture/internal/trainings/presentation/http"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rifat-simoom/go-clean-architecture/internal/shared_kernel/logs"
	"github.com/rifat-simoom/go-clean-architecture/internal/shared_kernel/server"
)

func main() {
	logs.Init()

	ctx := context.Background()

	app, cleanup := configs.NewApplication(ctx)
	defer cleanup()

	server.RunHTTPServer(func(router chi.Router) http.Handler {
		return http2.HandlerFromMux(http2.NewHttpServer(app), router)
	})
}
