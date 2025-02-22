package command

import (
	"context"
	hour2 "github.com/rifat-simoom/go-hexarch/internal/trainer/src/domain/hour"
	"time"

	"github.com/rifat-simoom/go-hexarch/internal/shared_kernel/decorator"
	"github.com/rifat-simoom/go-hexarch/internal/shared_kernel/errors"
	"github.com/sirupsen/logrus"
)

type MakeHoursUnavailable struct {
	Hours []time.Time
}

type MakeHoursUnavailableHandler decorator.CommandHandler[MakeHoursUnavailable]

type makeHoursUnavailableHandler struct {
	hourRepo hour2.Repository
}

func NewMakeHoursUnavailableHandler(
	hourRepo hour2.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) MakeHoursUnavailableHandler {
	if hourRepo == nil {
		panic("hourRepo is nil")
	}

	return decorator.ApplyCommandDecorators[MakeHoursUnavailable](
		makeHoursUnavailableHandler{hourRepo: hourRepo},
		logger,
		metricsClient,
	)
}

func (c makeHoursUnavailableHandler) Handle(ctx context.Context, cmd MakeHoursUnavailable) error {
	for _, hourToUpdate := range cmd.Hours {
		if err := c.hourRepo.UpdateHour(ctx, hourToUpdate, func(h *hour2.Hour) (*hour2.Hour, error) {
			if err := h.MakeNotAvailable(); err != nil {
				return nil, err
			}
			return h, nil
		}); err != nil {
			return errors.NewSlugError(err.Error(), "unable-to-update-availability")
		}
	}

	return nil
}
