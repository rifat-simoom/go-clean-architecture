package query

import (
	"context"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainer/src/application/interfaces/repositories"
	"time"

	"github.com/rifat-simoom/go-clean-architecture/internal/shared_kernel/decorator"
	"github.com/sirupsen/logrus"
)

type HourAvailability struct {
	Hour time.Time
}

type HourAvailabilityHandler decorator.QueryHandler[HourAvailability, bool]

type hourAvailabilityHandler struct {
	hourRepo repositories.Repository
}

func NewHourAvailabilityHandler(
	hourRepo repositories.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) HourAvailabilityHandler {
	if hourRepo == nil {
		panic("nil hourRepo")
	}

	return decorator.ApplyQueryDecorators[HourAvailability, bool](
		hourAvailabilityHandler{hourRepo: hourRepo},
		logger,
		metricsClient,
	)
}

func (h hourAvailabilityHandler) Handle(ctx context.Context, query HourAvailability) (bool, error) {
	hour, err := h.hourRepo.GetHour(ctx, query.Hour)
	if err != nil {
		return false, err
	}

	return hour.IsAvailable(), nil
}
