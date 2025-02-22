package command

import (
	"context"
	"github.com/rifat-simoom/go-hexarch/internal/trainings/src/application/interfaces/services"
	training2 "github.com/rifat-simoom/go-hexarch/internal/trainings/src/domain/training"

	"github.com/pkg/errors"
	"github.com/rifat-simoom/go-hexarch/internal/shared_kernel/decorator"
	"github.com/rifat-simoom/go-hexarch/internal/shared_kernel/logs"
	"github.com/sirupsen/logrus"
)

type CancelTraining struct {
	TrainingUUID string
	User         training2.User
}

type CancelTrainingHandler decorator.CommandHandler[CancelTraining]

type cancelTrainingHandler struct {
	repo           training2.Repository
	userService    services.UserService
	trainerService services.TrainerService
}

func NewCancelTrainingHandler(
	repo training2.Repository,
	userService services.UserService,
	trainerService services.TrainerService,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) decorator.CommandHandler[CancelTraining] {
	if repo == nil {
		panic("nil repo")
	}
	if userService == nil {
		panic("nil user service")
	}
	if trainerService == nil {
		panic("nil trainer service")
	}

	return decorator.ApplyCommandDecorators[CancelTraining](
		cancelTrainingHandler{repo: repo, userService: userService, trainerService: trainerService},
		logger,
		metricsClient,
	)
}

func (h cancelTrainingHandler) Handle(ctx context.Context, cmd CancelTraining) (err error) {
	defer func() {
		logs.LogCommandExecution("CancelTrainingHandler", cmd, err)
	}()

	return h.repo.UpdateTraining(
		ctx,
		cmd.TrainingUUID,
		cmd.User,
		func(ctx context.Context, tr *training2.Training) (*training2.Training, error) {
			if err := tr.Cancel(); err != nil {
				return nil, err
			}

			if balanceDelta := training2.CancelBalanceDelta(*tr, cmd.User.Type()); balanceDelta != 0 {
				err := h.userService.UpdateTrainingBalance(ctx, tr.UserUUID(), balanceDelta)
				if err != nil {
					return nil, errors.Wrap(err, "unable to change trainings balance")
				}
			}

			if err := h.trainerService.CancelTraining(ctx, tr.Time()); err != nil {
				return nil, errors.Wrap(err, "unable to cancel training")
			}

			return tr, nil
		},
	)
}
