package clickhouse

import (
	"time"
)

type RequestLog struct {
	Timestamp  time.Time
	Method     string
	Path       string
	Status     int
	DurationMs int
}
