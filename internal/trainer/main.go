package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rifat-simoom/go-clean-architecture/internal/common/genproto/trainer"
	"github.com/rifat-simoom/go-clean-architecture/internal/common/logs"
	"github.com/rifat-simoom/go-clean-architecture/internal/common/server"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainer/ports"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainer/service"
	"google.golang.org/grpc"
)

func main() {
	logs.Init()
	err := godotenv.Load("trainer.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	application := service.NewApplication(ctx)

	serverType := strings.ToLower(os.Getenv("SERVER_TO_RUN"))
	switch serverType {
	case "http":
		go loadFixtures(application)
		server.RunHTTPServer(func(router chi.Router) http.Handler {
			return ports.HandlerFromMux(
				ports.NewHttpServer(application),
				router,
			)
		})
	case "grpc":
		server.RunGRPCServer(func(server *grpc.Server) {
			svc := ports.NewGrpcServer(application)
			trainer.RegisterTrainerServiceServer(server, svc)
		})

	default:
		panic(fmt.Sprintf("server type '%s' is not supported", serverType))
	}
}
