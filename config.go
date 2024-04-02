package main

import (
	"os"
	"time"

	env "github.com/caarlos0/env/v10"
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
)

type Config struct {
	WeatherAPI string `env:"WEATHERAPI_KEY"`
	Model      string `env:"MODEL"`
	LogLevel   string `env:"LOG_LEVEL"`
	Logger     *log.Logger
}

func LoadConfig() *Config {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.TimeOnly,
		Level:           log.DebugLevel,
	})

	conf := &Config{
		Model:    "adrienbrault/gorilla-openfunctions-v2:Q4_K_M",
		LogLevel: "info",
		Logger:   logger,
	}

	if err := godotenv.Load(); err != nil {
		logger.Fatal("Error loading .env file")
	}
	if err := env.Parse(conf); err != nil {
		logger.Error(err)
	}

	conf.Logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		TimeFormat:      time.TimeOnly,
		Level:           log.ParseLevel(conf.LogLevel),
	})

	return conf
}
