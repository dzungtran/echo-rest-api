package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/dzungtran/echo-rest-api/modules/core/domains"
	"github.com/dzungtran/echo-rest-api/modules/core/dto"
	"github.com/dzungtran/echo-rest-api/pkg/constants"
	"github.com/labstack/echo/v4"
)

func (m *MiddlewareManager) FireBaseAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cli, err := m.appConf.FirebaseApp.Auth(c.Request().Context())
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
		}

		idToken := strings.Replace(c.Request().Header.Get("Authorization"), "Bearer ", "", 1)
		tkn, err := cli.VerifyIDToken(c.Request().Context(), idToken)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"error": err.Error(),
			})
		}

		tknUser := parseTokenClaimsToUser(tkn.Claims)
		var u *domains.UserWithRoles

		// Get user info
		u, err = m.fetchUserFromAuth(c.Request().Context(), tknUser.Code, tknUser.Email)
		if err != nil {
			if errors.Is(err, constants.ErrNotFound) {
				_, err = m.userUC.Register(c.Request().Context(), dto.CreateUserReq{
					FirstName: tknUser.FirstName,
					Email:     tknUser.Email,
					Code:      tknUser.Code,
				})

				if err != nil {
					return c.JSON(http.StatusInternalServerError, map[string]interface{}{
						"error": err.Error(),
					})
				}

				u, err = m.fetchUserFromAuth(c.Request().Context(), tknUser.Code, tknUser.Email)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, map[string]interface{}{
						"error": err.Error(),
					})
				}
			} else {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error": err.Error(),
				})
			}
		}

		if u.Status != domains.UserStatusActive {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"error": "user is not active",
			})
		}

		c.Set(constants.ContextKeyUser, u)
		return next(c)
	}
}

func parseTokenClaimsToUser(tknClaims map[string]interface{}) *domains.User {
	userExtId := tknClaims["user_id"].(string)
	userEmail := tknClaims["email"].(string)
	name := tknClaims["name"].(string)

	return &domains.User{
		Code:      userExtId,
		Email:     userEmail,
		FirstName: name,
	}
}
