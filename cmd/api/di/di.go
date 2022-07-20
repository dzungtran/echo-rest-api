package di

import (
	"github.com/dzungtran/echo-rest-api/config"
	"github.com/dzungtran/echo-rest-api/infrastructure/datastore"
	"github.com/dzungtran/echo-rest-api/modules/core"
	"github.com/dzungtran/echo-rest-api/modules/projects"
	"github.com/dzungtran/echo-rest-api/pkg/logger"
	"github.com/dzungtran/echo-rest-api/pkg/middlewares"
	sqlTools "github.com/dzungtran/echo-rest-api/pkg/sql-tools"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

func BuildDIContainer(
	mdbi *datastore.MasterDbInstance,
	sdbi *datastore.SlaveDbInstance,
	conf *config.AppConfig,
) *dig.Container {
	container := dig.New()
	_ = container.Provide(func() (*datastore.MasterDbInstance, *datastore.SlaveDbInstance) {
		return mdbi, sdbi
	})
	_ = container.Provide(func() *config.AppConfig {
		return conf
	})

	container.Provide(func() *sqlTools.SqlxTransaction {
		return sqlTools.NewSqlxTransaction(mdbi)
	})

	return container
}

func RegisterModules(e *echo.Echo, container *dig.Container) error {
	var err error
	mapModules := map[string]core.ModuleInstance{
		"core":     core.Module,
		"projects": projects.Module,
	}

	gRoot := e.Group("/")
	for _, m := range mapModules {
		err = m.RegisterRepositories(container)
		if err != nil {
			logger.Log().Errorf("RegisterRepositories error: %v", err)
			return err
		}

		err = m.RegisterUseCases(container)
		if err != nil {
			logger.Log().Errorf("RegisterUseCases error: %v", err)
			return err
		}
	}

	err = container.Provide(middlewares.NewMiddlewareManager)
	if err != nil {
		logger.Log().Errorf("RegisterHandlers error: %v", err)
		return err
	}

	for _, m := range mapModules {
		err = m.RegisterHandlers(gRoot, container)
		if err != nil {
			logger.Log().Errorf("RegisterHandlers error: %v", err)
			return err
		}
	}

	return err
}
