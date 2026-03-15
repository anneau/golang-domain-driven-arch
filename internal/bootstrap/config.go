package bootstrap

import (
	"os"
	"strconv"
	"strings"
)

type serverConfig struct {
	Port         int
	AllowOrigins []string
}

type dbConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
}

type auth0Config struct {
	Domain       string
	ClientID     string
	ClientSecret string
}

type sqsConfig struct {
	QueueURL          string
	EndpointURL       string
	Region            string
	MaxMessages       int
	WaitTimeSeconds   int
	VisibilityTimeout int
}

type config struct {
	AppEnv string
	Server serverConfig
	DB     dbConfig
	Auth0  auth0Config
	SQS    sqsConfig
}

func (c config) IsDevelopment() bool {
	return c.AppEnv == "" || strings.ToLower(c.AppEnv) == "development"
}

func (c config) IsProduction() bool {
	return strings.ToLower(c.AppEnv) == "production"
}

func loadConfig() config {
	return config{
		AppEnv: getEnv("APP_ENV", "development"),
		Server: serverConfig{
			Port:         getEnvInt("PORT", 8080),
			AllowOrigins: []string{getEnv("CORS_ALLOW_ORIGIN", "*")},
		},
		DB: dbConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "sample_db"),
		},
		Auth0: auth0Config{
			Domain:       getEnv("AUTH0_DOMAIN", ""),
			ClientID:     getEnv("AUTH0_CLIENT_ID", ""),
			ClientSecret: getEnv("AUTH0_CLIENT_SECRET", ""),
		},
		SQS: sqsConfig{
			QueueURL:          getEnv("SQS_QUEUE_URL", "http://localhost:4566/000000000000/events"),
			EndpointURL:       getEnv("SQS_ENDPOINT_URL", "http://localhost:4566"),
			Region:            getEnv("SQS_REGION", "ap-northeast-1"),
			MaxMessages:       getEnvInt("SQS_MAX_MESSAGES", 10),
			WaitTimeSeconds:   getEnvInt("SQS_WAIT_TIME_SECONDS", 20),
			VisibilityTimeout: getEnvInt("SQS_VISIBILITY_TIMEOUT", 30),
		},
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
