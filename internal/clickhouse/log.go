package clickhouse

import (
	"time"
)

type RequestLog struct {
	Timestamp  time.Time
	Method     string
	Path       string
	Status     uint16
	DurationMs uint32
}
