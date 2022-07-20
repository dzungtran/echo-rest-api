package core

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type ModuleInstance interface {
	RegisterRepositories(container *dig.Container) error
	RegisterUseCases(container *dig.Container) error
	RegisterHandlers(g *echo.Group, container *dig.Container) error
}
