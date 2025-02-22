package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainer/src/infrastructure/configs"
	presentation2 "github.com/rifat-simoom/go-clean-architecture/internal/trainer/src/presentation/grpc"
	http2 "github.com/rifat-simoom/go-clean-architecture/internal/trainer/src/presentation/http"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rifat-simoom/go-clean-architecture/internal/shared_kernel/genproto/trainer"
	"github.com/rifat-simoom/go-clean-architecture/internal/shared_kernel/logs"
	"github.com/rifat-simoom/go-clean-architecture/internal/shared_kernel/server"
	"google.golang.org/grpc"
)

func main() {
	logs.Init()
	err := godotenv.Load("trainer.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()

	application := configs.NewApplication(ctx)

	serverType := strings.ToLower(os.Getenv("SERVER_TO_RUN"))
	switch serverType {
	case "http":
		go configs.LoadFixtures(application)
		server.RunHTTPServer(func(router chi.Router) http.Handler {
			return http2.HandlerFromMux(
				http2.NewHttpServer(application),
				router,
			)
		})
	case "grpc":
		server.RunGRPCServer(func(server *grpc.Server) {
			svc := presentation2.NewGrpcServer(application)
			trainer.RegisterTrainerServiceServer(server, svc)
		})

	default:
		panic(fmt.Sprintf("server type '%s' is not supported", serverType))
	}
}
