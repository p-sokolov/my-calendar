package clickhouse

import (
	"context"
	"time"

	ch "github.com/ClickHouse/clickhouse-go/v2"
)

var Conn ch.Conn

func Connect() error {
	var err error

	Conn, err = ch.Open(&ch.Options{
		Addr: []string{"clickhouse:9000"},
		Auth: ch.Auth{
			Database: "calendar",
			Username: "calendar",
			Password: "calendar123",
		},
	})

	if err != nil {
		return err
	}

	return Conn.Ping(context.Background())
}

func ConnectWithRetry() error {
	var err error

	for i := 0; i < 10; i++ {
		err = Connect()
		if err == nil {
			return nil
		}

		time.Sleep(3 * time.Second)
	}

	return err
}

func SaveLog(log RequestLog) error {
	return Conn.Exec(
		context.Background(),
		`
        INSERT INTO request_logs
        (
            timestamp,
            method,
            path,
            status,
            duration_ms
        )
        VALUES (?, ?, ?, ?, ?)
        `,
		log.Timestamp,
		log.Method,
		log.Path,
		log.Status,
		log.DurationMs,
	)
}
