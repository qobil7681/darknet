package configs

import (
	"log"
	"os"
)

type LoggerConfig struct {
	Level   string
	Encoder string
}

func SetupLogger() {
	LoggerLevel := os.Getenv("LOGGER_LEVEL")
	if LoggerLevel == "" {
		log.Fatal("Log Level was not specified")
	}
	Encoder := os.Getenv("LOGGER_ENCODER")
	if Encoder == "" {
		log.Fatal("Logger encoder was not specified")
	}
	cfg.LoggerConf = LoggerConfig{
		Level:   LoggerLevel,
		Encoder: Encoder,
	}
}