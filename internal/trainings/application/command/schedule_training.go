package command

import (
	"context"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainings/application/interfaces/services"
	"time"

	"github.com/pkg/errors"
	"github.com/rifat-simoom/go-clean-architecture/internal/shared_kernel/decorator"
	"github.com/rifat-simoom/go-clean-architecture/internal/shared_kernel/logs"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainings/domain/training"
	"github.com/sirupsen/logrus"
)

type ScheduleTraining struct {
	TrainingUUID string

	UserUUID string
	UserName string

	TrainingTime time.Time
	Notes        string
}

type ScheduleTrainingHandler decorator.CommandHandler[ScheduleTraining]

type scheduleTrainingHandler struct {
	repo           training.Repository
	userService    services.UserService
	trainerService services.TrainerService
}

func NewScheduleTrainingHandler(
	repo training.Repository,
	userService services.UserService,
	trainerService services.TrainerService,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) ScheduleTrainingHandler {
	if repo == nil {
		panic("nil repo")
	}
	if userService == nil {
		panic("nil repo")
	}
	if trainerService == nil {
		panic("nil trainerService")
	}

	return decorator.ApplyCommandDecorators[ScheduleTraining](
		scheduleTrainingHandler{repo: repo, userService: userService, trainerService: trainerService},
		logger,
		metricsClient,
	)
}

func (h scheduleTrainingHandler) Handle(ctx context.Context, cmd ScheduleTraining) (err error) {
	defer func() {
		logs.LogCommandExecution("ScheduleTraining", cmd, err)
	}()

	tr, err := training.NewTraining(cmd.TrainingUUID, cmd.UserUUID, cmd.UserName, cmd.TrainingTime)
	if err != nil {
		return err
	}

	if err := h.repo.AddTraining(ctx, tr); err != nil {
		return err
	}

	err = h.userService.UpdateTrainingBalance(ctx, tr.UserUUID(), -1)
	if err != nil {
		return errors.Wrap(err, "unable to change trainings balance")
	}

	err = h.trainerService.ScheduleTraining(ctx, tr.Time())
	if err != nil {
		return errors.Wrap(err, "unable to schedule training")
	}

	return nil
}
