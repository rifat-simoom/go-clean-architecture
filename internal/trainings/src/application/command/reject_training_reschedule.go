package command

import (
	"context"
	training2 "github.com/rifat-simoom/go-clean-architecture/internal/trainings/src/domain/training"

	"github.com/rifat-simoom/go-clean-architecture/internal/shared_kernel/decorator"
	"github.com/rifat-simoom/go-clean-architecture/internal/shared_kernel/logs"
	"github.com/sirupsen/logrus"
)

type RejectTrainingReschedule struct {
	TrainingUUID string
	User         training2.User
}

type RejectTrainingRescheduleHandler decorator.CommandHandler[RejectTrainingReschedule]

type rejectTrainingRescheduleHandler struct {
	repo training2.Repository
}

func NewRejectTrainingRescheduleHandler(
	repo training2.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) RejectTrainingRescheduleHandler {
	if repo == nil {
		panic("nil repo service")
	}

	return decorator.ApplyCommandDecorators[RejectTrainingReschedule](
		rejectTrainingRescheduleHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h rejectTrainingRescheduleHandler) Handle(ctx context.Context, cmd RejectTrainingReschedule) (err error) {
	defer func() {
		logs.LogCommandExecution("RejectTrainingReschedule", cmd, err)
	}()

	return h.repo.UpdateTraining(
		ctx,
		cmd.TrainingUUID,
		cmd.User,
		func(ctx context.Context, tr *training2.Training) (*training2.Training, error) {
			if err := tr.RejectReschedule(); err != nil {
				return nil, err
			}

			return tr, nil
		},
	)
}
