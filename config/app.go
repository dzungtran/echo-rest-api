package config

import (
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// AppConfig - Init app config
type AppConfig struct {
	Environment string `json:"environment"`
	AppPort     string `json:"app_port"`
	DatabaseURL string `json:"database_url"`
	RedisURL    string `json:"redis_url"`

	Validator  echo.Validator        `json:"-"`
	CORSConfig middleware.CORSConfig `json:"-"`

	// 3rd-parties settings
	KratosWebhookApiKey string `json:"kratos_webhook_api_key"`
	KratosApiEndpoint   string `json:"kratos_api_endpoint"`
	AutoMigrate         bool   `json:"auto_migrate"`
	LogLevel            string `json:"log_level"`
}

type AppValidator struct {
	validator *validator.Validate
}

func (cv *AppValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func InitAppConfig() (*AppConfig, error) {
	currentEnv := "development"
	appPort := os.Getenv("PORT")
	if appPort == "" {
		appPort = "8088"
	}

	if os.Getenv("ENV") != "" {
		currentEnv = os.Getenv("ENV")
	}

	return &AppConfig{
		Environment: currentEnv,
		AppPort:     appPort,
		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisURL:    os.Getenv("REDIS_URL"),
		Validator:   &AppValidator{validator: validator.New()},
		CORSConfig:  middleware.DefaultCORSConfig,

		// 3rd-parties settings
		AutoMigrate:         os.Getenv("AUTO_MIGRATE") == "true",
		KratosWebhookApiKey: os.Getenv("KRATOS_WEBHOOK_API_KEY"),
		KratosApiEndpoint:   os.Getenv("KRATOS_API_ENDPOINT"),
		LogLevel:            os.Getenv("LOG_LEVEL"),
	}, nil
}
