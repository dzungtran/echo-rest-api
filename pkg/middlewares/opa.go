package middlewares

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/dzungtran/echo-rest-api/domains"
	"github.com/dzungtran/echo-rest-api/pkg/authz"
)

// Dểpcated: Call in handlers
func (m *MiddlewareManager) CheckPolicies(c echo.Context, callOpts ...authz.CallOPAInputOption) (denyMsg []string, err error) {
	u := c.Get("user")

	user, ok := u.(*domains.UserWithRoles)
	if !ok {
		err = errors.New("invalid user")
		return
	}

	callOpts = append(callOpts,
		authz.WithInputRequestMethod(c.Request().Method),
		authz.WithInputRequestEndpoint(c.Path()),
	)

	denyMsg, err = authz.CheckPolicies(user, callOpts...)
	if err == nil {
		c.Set("verified", true)
	}
	return
}
