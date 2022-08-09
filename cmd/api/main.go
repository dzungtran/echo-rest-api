package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/dzungtran/echo-rest-api/cmd/api/di"
	"github.com/dzungtran/echo-rest-api/config"
	_ "github.com/dzungtran/echo-rest-api/docs"
	"github.com/dzungtran/echo-rest-api/infrastructure/datastore"
	"github.com/dzungtran/echo-rest-api/migrations"
	"github.com/dzungtran/echo-rest-api/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Echo REST API
// @version 1.0
// @description This documentation for Echo REST server.
// @termsOfService http://swagger.io/terms/

// @contact.name Dzung Tran
// @contact.url https://docs.api.com/support
// @contact.email support@api.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host docs.api.com
// @BasePath /
func main() {
	// init app config
	conf, _ := config.InitAppConfig()

	logger.InitWithOptions(logger.WithConfigLevel(conf.LogLevel))
	if logger.Log() != nil {
		defer logger.Log().Sync()
	}

	// Echo instance
	e := echo.New()

	// Bind default middleware
	e.Use(middleware.LoggerWithConfig(config.GetEchoLogConfig(conf)))
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.HideBanner = true
	// e.HidePort = true
	e.Validator = conf.Validator

	// Setup infra
	mDBInstance := datastore.NewMasterDbInstance(conf.DatabaseURL)
	sDBInstance := datastore.NewSlaveDbInstance(conf.DatabaseURL)

	if conf.AutoMigrate {
		migrateDBInstance := datastore.NewMasterDbInstance(conf.DatabaseURL)
		migrations.RunAutoMigrate(migrateDBInstance.DBX().DB)
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

	err := di.RegisterModules(e, container)
	if err != nil {
		e.Logger.Fatal(err)
	}

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server
	go func() {
		if err := e.Start(":" + conf.AppPort); err != nil && err != http.ErrServerClosed {
			logger.Log().Fatal("shutting down the server")
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
