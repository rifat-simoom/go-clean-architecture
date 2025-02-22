package integration

import (
	"context"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainings/infrastructure/configs"
	http2 "github.com/rifat-simoom/go-clean-architecture/internal/trainings/presentation/http"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rifat-simoom/go-clean-architecture/internal/shared_kernel/server"
	"github.com/rifat-simoom/go-clean-architecture/internal/shared_kernel/tests"
	"github.com/stretchr/testify/require"
)

func TestCreateTraining(t *testing.T) {
	t.Parallel()

	token := tests.FakeAttendeeJWT(t, uuid.New().String())
	client := tests.NewTrainingsHTTPClient(t, token)

	hour := tests.RelativeDate(10, 12)
	trainingUUID := client.CreateTraining(t, "some note", hour)

	trainingsResponse := client.GetTrainings(t)

	var trainingsUUIDs []string
	for _, t := range trainingsResponse.Trainings {
		trainingsUUIDs = append(trainingsUUIDs, t.Uuid)
	}

	require.Contains(t, trainingsUUIDs, trainingUUID)
}

func TestCancelTraining(t *testing.T) {
	t.Parallel()

	token := tests.FakeAttendeeJWT(t, uuid.New().String())
	client := tests.NewTrainingsHTTPClient(t, token)

	hour := tests.RelativeDate(10, 13)
	trainingUUID := client.CreateTraining(t, "some note", hour)

	client.CancelTraining(t, trainingUUID, http.StatusOK)

	trainingsResponse := client.GetTrainings(t)

	var trainingsUUIDs []string
	for _, t := range trainingsResponse.Trainings {
		trainingsUUIDs = append(trainingsUUIDs, t.Uuid)
	}

	require.NotContains(t, trainingsUUIDs, trainingUUID)
}

func startService() bool {
	app := configs.NewComponentTestApplication(context.Background())

	trainingsHTTPAddr := os.Getenv("TRAININGS_HTTP_ADDR")
	go server.RunHTTPServerOnAddr(trainingsHTTPAddr, func(router chi.Router) http.Handler {
		return http2.HandlerFromMux(http2.NewHttpServer(app), router)
	})

	ok := tests.WaitForPort(trainingsHTTPAddr)
	if !ok {
		log.Println("Timed out waiting for trainings HTTP to come up")
	}

	return ok
}

func TestMain(m *testing.M) {
	if !startService() {
		os.Exit(1)
	}

	os.Exit(m.Run())
}
