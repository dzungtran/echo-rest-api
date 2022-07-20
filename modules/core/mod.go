package core

import (
	"github.com/dzungtran/echo-rest-api/modules/core/handlers"
	"github.com/dzungtran/echo-rest-api/modules/core/repositories"
	"github.com/dzungtran/echo-rest-api/modules/core/usecases"
	"github.com/dzungtran/echo-rest-api/pkg/middlewares"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

var Module ModuleInstance = &coreModule{}

type coreModule struct{}

func (coreModule) RegisterRepositories(container *dig.Container) error {
	container.Provide(repositories.NewPgsqlUserRepository)
	container.Provide(repositories.NewPgsqlOrgRepository)
	container.Provide(repositories.NewPgsqlUserOrgRepository)
	return nil
}

func (coreModule) RegisterUseCases(container *dig.Container) error {
	container.Provide(usecases.NewUserUsecase)
	container.Provide(usecases.NewOrgUsecase)
	return nil
}

func (coreModule) RegisterHandlers(g *echo.Group, container *dig.Container) error {
	return container.Invoke(func(
		middManager *middlewares.MiddlewareManager,
		userUsecase usecases.UserUsecase,
		orgUsecase usecases.OrgUsecase,
	) {
		handlers.NewOrgHandler(g, middManager, orgUsecase)
		handlers.NewUserHandler(g, middManager, userUsecase)
		handlers.NewKratosHookHandler(g, middManager, userUsecase)
	})
}
