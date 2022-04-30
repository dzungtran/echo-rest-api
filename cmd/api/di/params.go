package di

import (
	"github.com/dzungtran/echo-rest-api/config"
	"github.com/dzungtran/echo-rest-api/pkg/middlewares"
	"github.com/dzungtran/echo-rest-api/usecases"
	"go.uber.org/dig"
)

type DIContainerParams struct {
	dig.In
	AppConfig         *config.AppConfig
	MiddlewareManager *middlewares.MiddlewareManager
	UserUsecase       usecases.UserUsecase
	OrgUsecase        usecases.OrgUsecase
	// Auto generate
	// DO NOT DELETE THIS LINE ABOVE
}
