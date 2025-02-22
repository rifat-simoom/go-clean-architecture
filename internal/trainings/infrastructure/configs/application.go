package configs

import (
	"context"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainings/application"
	command2 "github.com/rifat-simoom/go-clean-architecture/internal/trainings/application/command"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainings/application/interfaces/services"
	query2 "github.com/rifat-simoom/go-clean-architecture/internal/trainings/application/query"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainings/infrastructure/persistence/respositories"
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
	trainerGrpc := respositories.NewTrainerGrpc(trainerClient)
	usersGrpc := respositories.NewUsersGrpc(usersClient)

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

	trainingsRepository := respositories.NewTrainingsFirestoreRepository(client)

	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NoOp{}

	return application.Application{
		Commands: application.Commands{
			ApproveTrainingReschedule: command2.NewApproveTrainingRescheduleHandler(trainingsRepository, usersGrpc, trainerGrpc, logger, metricsClient),
			CancelTraining:            command2.NewCancelTrainingHandler(trainingsRepository, usersGrpc, trainerGrpc, logger, metricsClient),
			RejectTrainingReschedule:  command2.NewRejectTrainingRescheduleHandler(trainingsRepository, logger, metricsClient),
			RescheduleTraining:        command2.NewRescheduleTrainingHandler(trainingsRepository, usersGrpc, trainerGrpc, logger, metricsClient),
			RequestTrainingReschedule: command2.NewRequestTrainingRescheduleHandler(trainingsRepository, logger, metricsClient),
			ScheduleTraining:          command2.NewScheduleTrainingHandler(trainingsRepository, usersGrpc, trainerGrpc, logger, metricsClient),
		},
		Queries: application.Queries{
			AllTrainings:     query2.NewAllTrainingsHandler(trainingsRepository, logger, metricsClient),
			TrainingsForUser: query2.NewTrainingsForUserHandler(trainingsRepository, logger, metricsClient),
		},
	}
}
