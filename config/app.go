package config

import (
	"context"
	"fmt"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"google.golang.org/api/option"
)

// AppConfig - Init app config
type AppConfig struct {
	Environment string `json:"environment"`
	AppPort     string `json:"app_port"`
	BaseURL     string `json:"base_url"`
	DatabaseURL string `json:"database_url"`
	RedisURL    string `json:"redis_url"`

	Validator   echo.Validator        `json:"-"`
	CORSConfig  middleware.CORSConfig `json:"-"`
	FirebaseApp *firebase.App         `json:"-"`

	// 3rd-parties settings
	AutoMigrate bool   `json:"auto_migrate"`
	LogLevel    string `json:"log_level"`

	AuthProvider        string `json:"auth_provider"`
	FirebaseCreds       string `json:"firebase_creds"`
	FirebaseAuthCreds   string `json:"firebase_auth_creds"`
	KratosWebhookApiKey string `json:"kratos_webhook_api_key"`
	KratosApiEndpoint   string `json:"kratos_api_endpoint"`
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
	var fbApp *firebase.App
	var err error
	currentEnv := "development"

	// load env file if exists
	_, err = os.Stat(".env")
	if err == nil {
		err = godotenv.Load(os.ExpandEnv(".env"))
		if err != nil {
			return nil, fmt.Errorf("error initializing app: %v", err)
		}
	}

	appPort := os.Getenv("PORT")
	if appPort == "" {
		appPort = "8088"
	}

	if os.Getenv("ENV") != "" {
		currentEnv = os.Getenv("ENV")
	}

	fbOpt := option.WithCredentialsJSON([]byte(os.Getenv("FIREBASE_CREDENTIALS")))
	if os.Getenv("FIREBASE_CREDENTIALS") != "" {
		fbApp, err = firebase.NewApp(context.Background(), nil, fbOpt)
		if err != nil {
			return nil, fmt.Errorf("error initializing app: %v", err)
		}
	}

	return &AppConfig{
		Environment: currentEnv,
		AppPort:     appPort,
		BaseURL:     os.Getenv("BASE_URL"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisURL:    os.Getenv("REDIS_URL"),
		Validator:   &AppValidator{validator: validator.New()},
		CORSConfig:  middleware.DefaultCORSConfig,

		// 3rd-parties settings
		AutoMigrate: os.Getenv("AUTO_MIGRATE") == "true",
		LogLevel:    os.Getenv("LOG_LEVEL"),

		AuthProvider:        os.Getenv("AUTH_PROVIDER"),
		KratosWebhookApiKey: os.Getenv("KRATOS_WEBHOOK_API_KEY"),
		KratosApiEndpoint:   os.Getenv("KRATOS_API_ENDPOINT"),
		FirebaseApp:         fbApp,
		FirebaseCreds:       os.Getenv("FIREBASE_CREDENTIALS"),
		FirebaseAuthCreds:   os.Getenv("FIREBASE_AUTH_CREDENTIALS"),
	}, nil
}
