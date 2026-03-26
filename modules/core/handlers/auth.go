package handlers

import (
	"net/http"

	"github.com/dzungtran/echo-rest-api/config"
	"github.com/dzungtran/echo-rest-api/modules/core/usecases"
	"github.com/dzungtran/echo-rest-api/pkg/middlewares"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	UserUC  usecases.UserUsecase
	Configs *config.AppConfig
}

func NewAuthHandler(g *echo.Group, middManager *middlewares.MiddlewareManager, userUsecase usecases.UserUsecase, appCfgs *config.AppConfig) {
	handler := &AuthHandler{
		UserUC:  userUsecase,
		Configs: appCfgs,
	}

	apiV1 := g.Group("auth")
	apiV1.GET("/login", handler.LoginForm)
	apiV1.GET("/success", handler.LoginSuccess)
}

// LoginForm godoc
// @Summary      Show login form
// @Description  Renders Firebase login UI
// @Tags         auth
// @Produce      html
// @Success      200  {string}  string  "HTML login page"
// @Router       /auth/login [get]
func (h *AuthHandler) LoginForm(c echo.Context) error {
	data := map[string]interface{}{
		"config":      h.Configs.FirebaseAuthCreds,
		"success_url": h.Configs.BaseURL + "auth/success",
	}
	return c.Render(http.StatusOK, "login.go.tpl", data)
}

// LoginSuccess godoc
// @Summary      Show login success
// @Description  Renders Firebase login success UI
// @Tags         auth
// @Produce      html
// @Success      200  {string}  string  "HTML success page"
// @Router       /auth/success [get]
func (h *AuthHandler) LoginSuccess(c echo.Context) error {
	data := map[string]interface{}{
		"config": h.Configs.FirebaseAuthCreds,
	}
	return c.Render(http.StatusOK, "success.go.tpl", data)
}
