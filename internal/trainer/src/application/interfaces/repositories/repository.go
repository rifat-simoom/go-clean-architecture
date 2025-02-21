package repositories

import (
	"context"
	"github.com/rifat-simoom/go-clean-architecture/internal/trainer/src/domain/hour"
	"time"
)

type Repository interface {
	GetHour(ctx context.Context, hourTime time.Time) (*hour.Hour, error)
	UpdateHour(
		ctx context.Context,
		hourTime time.Time,
		updateFn func(h *hour.Hour) (*hour.Hour, error),
	) error
}
