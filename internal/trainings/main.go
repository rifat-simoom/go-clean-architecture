package main

import (
	"context"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainings/src/infrastructure/configs"
	http3 "github.com/rifat-simoom/go-clean-architecture/internal/trainings/src/presentation/http"
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
		return http3.HandlerFromMux(http3.NewHttpServer(app), router)
	})
}
