package middlewares

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RequireResourceIdInParam(paramName string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id, err := strconv.Atoi(c.Param(paramName))
			if err != nil || id <= 0 {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"message": paramName + " is invalid",
				})
			}
			return next(c)
		}
	}
}
