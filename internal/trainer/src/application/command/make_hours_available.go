package command

import (
	"context"
	hour2 "github.com/rifat-simoom/go-hexarch/internal/trainer/src/domain/hour"
	"time"

	"github.com/rifat-simoom/go-hexarch/internal/shared_kernel/decorator"
	"github.com/rifat-simoom/go-hexarch/internal/shared_kernel/errors"
	"github.com/sirupsen/logrus"
)

type MakeHoursAvailable struct {
	Hours []time.Time
}

type MakeHoursAvailableHandler decorator.CommandHandler[MakeHoursAvailable]

type makeHoursAvailableHandler struct {
	hourRepo hour2.Repository
}

func NewMakeHoursAvailableHandler(
	hourRepo hour2.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) MakeHoursAvailableHandler {
	if hourRepo == nil {
		panic("hourRepo is nil")
	}

	return decorator.ApplyCommandDecorators[MakeHoursAvailable](
		makeHoursAvailableHandler{hourRepo: hourRepo},
		logger,
		metricsClient,
	)
}

func (c makeHoursAvailableHandler) Handle(ctx context.Context, cmd MakeHoursAvailable) error {
	for _, hourToUpdate := range cmd.Hours {
		if err := c.hourRepo.UpdateHour(ctx, hourToUpdate, func(h *hour2.Hour) (*hour2.Hour, error) {
			if err := h.MakeAvailable(); err != nil {
				return nil, err
			}
			return h, nil
		}); err != nil {
			return errors.NewSlugError(err.Error(), "unable-to-update-availability")
		}
	}

	return nil
}
