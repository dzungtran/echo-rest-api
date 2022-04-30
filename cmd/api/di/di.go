package di

import (
	"github.com/dzungtran/echo-rest-api/config"
	httpDelivery "github.com/dzungtran/echo-rest-api/delivery/http"
	"github.com/dzungtran/echo-rest-api/infrastructure/datastore"
	"github.com/dzungtran/echo-rest-api/pkg/middlewares"
	"github.com/dzungtran/echo-rest-api/repositories/postgres"
	"github.com/dzungtran/echo-rest-api/usecases"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	_ = postgres.Inject(container)
	_ = usecases.Inject(container)
	_ = container.Provide(middlewares.NewMiddlewareManager)

	return container
}

func RegisterHandlers(e *echo.Echo, container *dig.Container) error {
	return container.Invoke(func(params DIContainerParams) {
		e.Use(params.MiddlewareManager.GenerateRequestID())

		// Group routes for middlewares
		adminGroup := e.Group("/admin",
			middleware.CORSWithConfig(params.AppConfig.CORSConfig),
		)

		hookGroup := e.Group("/hooks",
			middleware.CORSWithConfig(params.AppConfig.CORSConfig),
		)

		// bind api handlers to group
		httpDelivery.NewUserHandler(adminGroup, params.MiddlewareManager, params.UserUsecase)
		httpDelivery.NewOrgHandler(adminGroup, params.MiddlewareManager, params.OrgUsecase)
		httpDelivery.NewKratosHookHandler(hookGroup, params.MiddlewareManager, params.UserUsecase)
		// Auto generate
		// DO NOT DELETE THIS LINE ABOVE
	})
}
