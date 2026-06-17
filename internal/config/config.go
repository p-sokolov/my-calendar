package config

import (
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"my-calendar/internal/logger"
)

type Config struct {
	HttpPort string
}

func LoadCfg() *Config {
	err := godotenv.Load()
	if err != nil {
		logger.L().Fatal("config load failed", zap.Error(err))
	}

	httpPort := os.Getenv("HTTP_PORT")

	return &Config{HttpPort: httpPort}
}
