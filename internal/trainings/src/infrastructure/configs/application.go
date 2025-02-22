package configs

import (
	"context"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainings/src/application"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainings/src/application/command"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainings/src/application/interfaces/services"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainings/src/application/query"
	respositories2 "github.com/rifat-simoom/go-clean-architecture/internal/trainings/src/infrastructure/persistence/respositories"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainings/tests/integration"
	"os"

	"cloud.google.com/go/firestore"
	grpcClient "github.com/rifat-simoom/go-clean-architecture/internal/shared_kernel/client"
	"github.com/rifat-simoom/go-clean-architecture/internal/shared_kernel/metrics"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) (application.Application, func()) {
	trainerClient, closeTrainerClient, err := grpcClient.NewTrainerClient()
	if err != nil {
		panic(err)
	}

	usersClient, closeUsersClient, err := grpcClient.NewUsersClient()
	if err != nil {
		panic(err)
	}
	trainerGrpc := respositories2.NewTrainerGrpc(trainerClient)
	usersGrpc := respositories2.NewUsersGrpc(usersClient)

	return newApplication(ctx, trainerGrpc, usersGrpc),
		func() {
			_ = closeTrainerClient()
			_ = closeUsersClient()
		}
}

func NewComponentTestApplication(ctx context.Context) application.Application {
	return newApplication(ctx, integration.TrainerServiceMock{}, integration.UserServiceMock{})
}

func newApplication(ctx context.Context, trainerGrpc services.TrainerService, usersGrpc services.UserService) application.Application {
	client, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	if err != nil {
		panic(err)
	}

	trainingsRepository := respositories2.NewTrainingsFirestoreRepository(client)

	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NoOp{}

	return application.Application{
		Commands: application.Commands{
			ApproveTrainingReschedule: command.NewApproveTrainingRescheduleHandler(trainingsRepository, usersGrpc, trainerGrpc, logger, metricsClient),
			CancelTraining:            command.NewCancelTrainingHandler(trainingsRepository, usersGrpc, trainerGrpc, logger, metricsClient),
			RejectTrainingReschedule:  command.NewRejectTrainingRescheduleHandler(trainingsRepository, logger, metricsClient),
			RescheduleTraining:        command.NewRescheduleTrainingHandler(trainingsRepository, usersGrpc, trainerGrpc, logger, metricsClient),
			RequestTrainingReschedule: command.NewRequestTrainingRescheduleHandler(trainingsRepository, logger, metricsClient),
			ScheduleTraining:          command.NewScheduleTrainingHandler(trainingsRepository, usersGrpc, trainerGrpc, logger, metricsClient),
		},
		Queries: application.Queries{
			AllTrainings:     query.NewAllTrainingsHandler(trainingsRepository, logger, metricsClient),
			TrainingsForUser: query.NewTrainingsForUserHandler(trainingsRepository, logger, metricsClient),
		},
	}
}
