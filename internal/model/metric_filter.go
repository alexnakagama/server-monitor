package model

import "time"

type MetricFilters struct {
	From  *time.Time
	To    *time.Time
	Limit int
}
