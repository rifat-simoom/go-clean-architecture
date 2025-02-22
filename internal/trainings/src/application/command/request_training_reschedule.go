package command

import (
	"context"
	training2 "github.com/rifat-simoom/go-clean-architecture/internal/trainings/src/domain/training"
	"time"

	"github.com/rifat-simoom/go-clean-architecture/internal/shared_kernel/decorator"
	"github.com/rifat-simoom/go-clean-architecture/internal/shared_kernel/logs"
	"github.com/sirupsen/logrus"
)

type RequestTrainingReschedule struct {
	TrainingUUID string
	NewTime      time.Time

	User training2.User

	NewNotes string
}

type RequestTrainingRescheduleHandler decorator.CommandHandler[RequestTrainingReschedule]

type requestTrainingRescheduleHandler struct {
	repo training2.Repository
}

func NewRequestTrainingRescheduleHandler(
	repo training2.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) RequestTrainingRescheduleHandler {
	if repo == nil {
		panic("nil repo service")
	}

	return decorator.ApplyCommandDecorators[RequestTrainingReschedule](
		requestTrainingRescheduleHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h requestTrainingRescheduleHandler) Handle(ctx context.Context, cmd RequestTrainingReschedule) (err error) {
	defer func() {
		logs.LogCommandExecution("RequestTrainingReschedule", cmd, err)
	}()

	return h.repo.UpdateTraining(
		ctx,
		cmd.TrainingUUID,
		cmd.User,
		func(ctx context.Context, tr *training2.Training) (*training2.Training, error) {
			if err := tr.UpdateNotes(cmd.NewNotes); err != nil {
				return nil, err
			}

			tr.ProposeReschedule(cmd.NewTime, cmd.User.Type())

			return tr, nil
		},
	)
}
