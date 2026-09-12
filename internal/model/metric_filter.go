package model

import (
	"time"

	"github.com/alexnakagama/server-monitor/internal/server/errors_custom"
)

type MetricFilters struct {
	From  *time.Time
	To    *time.Time
	Limit int
}

func (f MetricFilters) Validate() error {
	if f.Limit < 1 || f.Limit > 1000 {
		return errors_custom.ErrInvalidLimit
	}

	if f.From != nil && f.To != nil && f.From.After(*f.To) {
		return errors_custom.ErrInvalidTimeRange
	}

	return nil
}
