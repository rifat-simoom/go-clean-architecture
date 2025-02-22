package command

import (
	"context"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainings/src/application/interfaces/services"
	training2 "github.com/rifat-simoom/go-clean-architecture/internal/trainings/src/domain/training"

	"github.com/rifat-simoom/go-clean-architecture/internal/shared_kernel/decorator"
	"github.com/rifat-simoom/go-clean-architecture/internal/shared_kernel/logs"
	"github.com/sirupsen/logrus"
)

type ApproveTrainingReschedule struct {
	TrainingUUID string
	User         training2.User
}

type ApproveTrainingRescheduleHandler decorator.CommandHandler[ApproveTrainingReschedule]

type approveTrainingRescheduleHandler struct {
	repo           training2.Repository
	userService    services.UserService
	trainerService services.TrainerService
}

func NewApproveTrainingRescheduleHandler(
	repo training2.Repository,
	userService services.UserService,
	trainerService services.TrainerService,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) decorator.CommandHandler[ApproveTrainingReschedule] {
	if repo == nil {
		panic("nil repo")
	}
	if userService == nil {
		panic("nil userService")
	}
	if trainerService == nil {
		panic("nil trainerService")
	}

	return decorator.ApplyCommandDecorators[ApproveTrainingReschedule](
		approveTrainingRescheduleHandler{repo, userService, trainerService},
		logger,
		metricsClient,
	)
}

func (h approveTrainingRescheduleHandler) Handle(ctx context.Context, cmd ApproveTrainingReschedule) (err error) {
	defer func() {
		logs.LogCommandExecution("ApproveTrainingReschedule", cmd, err)
	}()

	return h.repo.UpdateTraining(
		ctx,
		cmd.TrainingUUID,
		cmd.User,
		func(ctx context.Context, tr *training2.Training) (*training2.Training, error) {
			originalTrainingTime := tr.Time()

			if err := tr.ApproveReschedule(cmd.User.Type()); err != nil {
				return nil, err
			}

			err := h.trainerService.MoveTraining(ctx, tr.Time(), originalTrainingTime)
			if err != nil {
				return nil, err
			}

			return tr, nil
		},
	)
}
