package app

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	AppName                = "httpmux"
	UserAgent              = "sparfenov/" + AppName
	ListenLimit            = 100
	RequestURLLimit        = 20
	ExternalRequestLimit   = 4
	ExternalRequestTimeout = 1 * time.Second
	ServerShutdownTimeout  = 30 * time.Second
)

type Config struct {
	IsDebug                bool
	UserAgent              string
	Addr                   string
	ListenLimit            uint
	MaxURLCountToProcess   int
	ExternalRequestLimit   int
	ExternalRequestTimeout time.Duration
	ServerShutdownTimeout  time.Duration
}

func NewConfig() (*Config, error) {
	isDebug, err := getBoolEnv("IS_DEBUG", false)
	if err != nil {
		return nil, err
	}

	config := Config{
		IsDebug:                isDebug,
		UserAgent:              UserAgent,
		Addr:                   getStringEnv("SERVER_ADDR", ":8080"),
		ListenLimit:            ListenLimit,
		MaxURLCountToProcess:   RequestURLLimit,
		ExternalRequestLimit:   ExternalRequestLimit,
		ExternalRequestTimeout: ExternalRequestTimeout,
		ServerShutdownTimeout:  ServerShutdownTimeout,
	}

	return &config, nil
}

func getStringEnv(name string, defaultVal string) string {
	val := os.Getenv(name)
	if strings.TrimSpace(val) == "" {
		return defaultVal
	}

	return val
}

func getBoolEnv(name string, defaultVal bool) (bool, error) {
	valStr := os.Getenv(name)
	if valStr == "" {
		return defaultVal, nil
	}

	value, err := strconv.ParseBool(valStr)
	if err != nil {
		return false, fmt.Errorf("failed to parse bool for env variable: %s, error: %w", name, err)
	}

	return value, nil
}
