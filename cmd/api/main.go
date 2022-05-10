package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dzungtran/echo-rest-api/cmd/api/di"
	"github.com/dzungtran/echo-rest-api/config"
	"github.com/dzungtran/echo-rest-api/infrastructure/datastore"
	"github.com/dzungtran/echo-rest-api/migrations"
	"github.com/dzungtran/echo-rest-api/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// init app config
	conf, _ := config.InitAppConfig()

	logger.InitLog(conf.Environment)
	if logger.Log() != nil {
		defer logger.Log().Sync()
	}

	// Echo instance
	e := echo.New()

	// Bind default middleware
	e.Use(middleware.LoggerWithConfig(config.GetEchoLogConfig(conf)))
	e.Use(middleware.Recover())
	e.Validator = conf.Validator

	// Setup infra
	mDBInstance := datastore.NewMasterDbInstance(conf.DatabaseURL)
	sDBInstance := datastore.NewSlaveDbInstance(conf.DatabaseURL)
	if conf.AutoMigrate {
		migrations.RunAutoMigrate(mDBInstance.DBX().DB)
	}

	// Setup middleware manager
	defer func() {
		mDBInstance.DBX().Close()
		sDBInstance.DBX().Close()
	}()

	container := di.BuildDIContainer(
		mDBInstance,
		sDBInstance,
		conf,
	)

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Hello there!",
		})
	})

	err := di.RegisterHandlers(e, container)
	if err != nil {
		e.Logger.Fatal(err)
	}

	// Start server
	go func() {
		if err := e.Start(":" + conf.AppPort); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
