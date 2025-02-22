package query

import (
	"context"
	"time"

	"github.com/rifat-simoom/go-hexarch/internal/shared_kernel/decorator"
	"github.com/rifat-simoom/go-hexarch/internal/shared_kernel/errors"
	"github.com/sirupsen/logrus"
)

type AvailableHours struct {
	From time.Time
	To   time.Time
}

type AvailableHoursHandler decorator.QueryHandler[AvailableHours, []Date]

type AvailableHoursReadModel interface {
	AvailableHours(ctx context.Context, from time.Time, to time.Time) ([]Date, error)
}

type availableHoursHandler struct {
	readModel AvailableHoursReadModel
}

func NewAvailableHoursHandler(
	readModel AvailableHoursReadModel,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) AvailableHoursHandler {
	return decorator.ApplyQueryDecorators[AvailableHours, []Date](
		availableHoursHandler{readModel: readModel},
		logger,
		metricsClient,
	)
}

func (h availableHoursHandler) Handle(ctx context.Context, query AvailableHours) (d []Date, err error) {
	if query.From.After(query.To) {
		return nil, errors.NewIncorrectInputError("date-from-after-date-to", "Date from after date to")
	}

	return h.readModel.AvailableHours(ctx, query.From, query.To)
}
