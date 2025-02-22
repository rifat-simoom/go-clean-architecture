package configs

import (
	"context"
	"github.com/rifat-simoom/go-hexarch/internal/trainer/src/application"
	"github.com/rifat-simoom/go-hexarch/internal/trainer/src/application/command"
	"github.com/rifat-simoom/go-hexarch/internal/trainer/src/application/query"
	"github.com/rifat-simoom/go-hexarch/internal/trainer/src/domain/hour"
	repositories2 "github.com/rifat-simoom/go-hexarch/internal/trainer/src/infrastructure/persistence/repositories"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/rifat-simoom/go-hexarch/internal/shared_kernel/metrics"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) application.Application {
	firestoreClient, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	if err != nil {
		panic(err)
	}

	factoryConfig := hour.FactoryConfig{
		MaxWeeksInTheFutureToSet: 6,
		MinUtcHour:               12,
		MaxUtcHour:               20,
	}

	datesRepository := repositories2.NewDatesFirestoreRepository(firestoreClient, factoryConfig)

	hourFactory, err := hour.NewFactory(factoryConfig)
	if err != nil {
		panic(err)
	}

	hourRepository := repositories2.NewFirestoreHourRepository(firestoreClient, hourFactory)

	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NoOp{}

	return application.Application{
		Commands: application.Commands{
			CancelTraining:       command.NewCancelTrainingHandler(hourRepository, logger, metricsClient),
			ScheduleTraining:     command.NewScheduleTrainingHandler(hourRepository, logger, metricsClient),
			MakeHoursAvailable:   command.NewMakeHoursAvailableHandler(hourRepository, logger, metricsClient),
			MakeHoursUnavailable: command.NewMakeHoursUnavailableHandler(hourRepository, logger, metricsClient),
		},
		Queries: application.Queries{
			HourAvailability:      query.NewHourAvailabilityHandler(hourRepository, logger, metricsClient),
			TrainerAvailableHours: query.NewAvailableHoursHandler(datesRepository, logger, metricsClient),
		},
	}
}
