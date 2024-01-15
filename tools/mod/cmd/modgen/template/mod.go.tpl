package {{ .PluralName | ToKebab }}

import (
	"{{ .RootPackage }}/modules/core"
	"{{ .RootPackage }}/modules/{{ .PluralName | ToKebab }}/handlers"
	"{{ .RootPackage }}/modules/{{ .PluralName | ToKebab }}/repositories"
	"{{ .RootPackage }}/modules/{{ .PluralName | ToKebab }}/usecases"
	"{{ .RootPackage }}/pkg/middlewares"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

var Module core.ModuleInstance = &{{ .SingularName | ToLowerCamel }}Module{}

type {{ .SingularName | ToLowerCamel }}Module struct{}

func ({{ .SingularName | ToLowerCamel }}Module) RegisterRepositories(container *dig.Container) error {
	container.Provide(repositories.NewPgsql{{ .SingularName }}Repository)
	return nil
}

func ({{ .SingularName | ToLowerCamel }}Module) RegisterUseCases(container *dig.Container) error {
	container.Provide(usecases.New{{ .SingularName }}Usecase)
	return nil
}

func ({{ .SingularName | ToLowerCamel }}Module) RegisterHandlers(g *echo.Group, container *dig.Container) error {
	return container.Invoke(func(
		middManager *middlewares.MiddlewareManager,
		{{ .SingularName | ToLowerCamel }}Usecase usecases.{{ .SingularName }}Usecase,
	) {
		handlers.New{{ .SingularName }}Handler(g, middManager, {{ .SingularName | ToLowerCamel }}Usecase)
	})
}
