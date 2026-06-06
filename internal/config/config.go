package config

import (
	"os"

	"github.com/joho/godotenv"
	_"go.uber.org/zap"
)

type Config struct {
	HttpPort string
}

func LoadCfg() *Config {
	err := godotenv.Load()
	if err != nil {
		// logger.fatal
	}

	httpPort := os.Getenv("HTTP_PORT")

	return &Config{ HttpPort: httpPort }
}