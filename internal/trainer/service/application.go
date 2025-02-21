package service

import (
	"context"
	"os"

	"cloud.google.com/go/firestore"
	"github.com/rifat-simoom/go-clean-architecture/internal/shared_kernel/metrics"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainer/adapters"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainer/app"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainer/app/command"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainer/app/query"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainer/domain/hour"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context) app.Application {
	firestoreClient, err := firestore.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	if err != nil {
		panic(err)
	}

	factoryConfig := hour.FactoryConfig{
		MaxWeeksInTheFutureToSet: 6,
		MinUtcHour:               12,
		MaxUtcHour:               20,
	}

	datesRepository := adapters.NewDatesFirestoreRepository(firestoreClient, factoryConfig)

	hourFactory, err := hour.NewFactory(factoryConfig)
	if err != nil {
		panic(err)
	}

	hourRepository := adapters.NewFirestoreHourRepository(firestoreClient, hourFactory)

	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NoOp{}

	return app.Application{
		Commands: app.Commands{
			CancelTraining:       command.NewCancelTrainingHandler(hourRepository, logger, metricsClient),
			ScheduleTraining:     command.NewScheduleTrainingHandler(hourRepository, logger, metricsClient),
			MakeHoursAvailable:   command.NewMakeHoursAvailableHandler(hourRepository, logger, metricsClient),
			MakeHoursUnavailable: command.NewMakeHoursUnavailableHandler(hourRepository, logger, metricsClient),
		},
		Queries: app.Queries{
			HourAvailability:      query.NewHourAvailabilityHandler(hourRepository, logger, metricsClient),
			TrainerAvailableHours: query.NewAvailableHoursHandler(datesRepository, logger, metricsClient),
		},
	}
}
