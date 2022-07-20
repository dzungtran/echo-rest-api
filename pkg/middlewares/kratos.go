package middlewares

import (
	"net/http"

	"github.com/dzungtran/echo-rest-api/modules/core/domains"
	"github.com/dzungtran/echo-rest-api/pkg/constants"
	"github.com/dzungtran/echo-rest-api/pkg/logger"
	"github.com/labstack/echo/v4"
	ory "github.com/ory/kratos-client-go"
)

func (m *MiddlewareManager) KratosWebhookAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if m.appConf.KratosWebhookApiKey != c.Request().Header.Get(constants.HeaderXApiKey) {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "Unauthorized",
			})
		}
		return next(c)
	}
}

func (m *MiddlewareManager) KratosAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		var ss *ory.Session
		var err error

		ctx := c.Request().Context()
		ss, _, err = m.kratosClient.V0alpha2Api.ToSession(ctx).
			Cookie(c.Request().Header.Get("Cookie")).
			XSessionToken(c.Request().Header.Get("X-Session-Token")).
			Execute()

		if err != nil {
			logger.Log().Errorw("error while fetch identity from kratos", "error", err)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": err.Error(),
			})
		}

		if ss.Active == nil || !*ss.Active {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "user is inactive",
			})
		}

		_, ok := ss.Identity.GetTraitsOk()
		if !ok {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "cannot get identity info",
			})
		}

		u, err := m.fetchUserFromAuth(ctx, ss.Identity.Id, "")
		if err != nil {
			logger.Log().Errorw("error while fetch user", "error", err, "code", ss.Identity.Id)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "cannot fetch user info",
			})
		}

		if u.Status != domains.UserStatusActive {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"message": "user is not active",
			})
		}

		c.Set(constants.ContextKeyUser, u)
		return next(c)
	}
}
