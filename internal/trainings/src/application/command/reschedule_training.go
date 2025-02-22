package command

import (
	"context"
	"github.com/rifat-simoom/go-hexarch/internal/trainings/src/application/interfaces/services"
	training2 "github.com/rifat-simoom/go-hexarch/internal/trainings/src/domain/training"
	"time"

	"github.com/rifat-simoom/go-hexarch/internal/shared_kernel/decorator"
	"github.com/rifat-simoom/go-hexarch/internal/shared_kernel/logs"
	"github.com/sirupsen/logrus"
)

type RescheduleTraining struct {
	TrainingUUID string
	NewTime      time.Time

	User training2.User

	NewNotes string
}

type RescheduleTrainingHandler decorator.CommandHandler[RescheduleTraining]

type rescheduleTrainingHandler struct {
	repo           training2.Repository
	userService    services.UserService
	trainerService services.TrainerService
}

func NewRescheduleTrainingHandler(
	repo training2.Repository,
	userService services.UserService,
	trainerService services.TrainerService,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) RescheduleTrainingHandler {
	if repo == nil {
		panic("nil repo")
	}
	if userService == nil {
		panic("nil userService")
	}
	if trainerService == nil {
		panic("nil trainerService")
	}

	return decorator.ApplyCommandDecorators[RescheduleTraining](
		rescheduleTrainingHandler{repo: repo, userService: userService, trainerService: trainerService},
		logger,
		metricsClient,
	)
}

func (h rescheduleTrainingHandler) Handle(ctx context.Context, cmd RescheduleTraining) (err error) {
	defer func() {
		logs.LogCommandExecution("RescheduleTraining", cmd, err)
	}()

	return h.repo.UpdateTraining(
		ctx,
		cmd.TrainingUUID,
		cmd.User,
		func(ctx context.Context, tr *training2.Training) (*training2.Training, error) {
			originalTrainingTime := tr.Time()

			if err := tr.UpdateNotes(cmd.NewNotes); err != nil {
				return nil, err
			}

			if err := tr.RescheduleTraining(cmd.NewTime); err != nil {
				return nil, err
			}

			err := h.trainerService.MoveTraining(ctx, cmd.NewTime, originalTrainingTime)
			if err != nil {
				return nil, err
			}

			return tr, nil
		},
	)
}
