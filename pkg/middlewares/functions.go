package middlewares

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/dzungtran/echo-rest-api/pkg/constants"
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

func BindRequestPayload(payloadInst interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if reflect.ValueOf(payloadInst).Kind() != reflect.Ptr {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error": "payload should be a pointer",
				})
			}

			if err := c.Bind(payloadInst); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"error": err.Error(),
				})
			}

			c.Set(constants.ContextKeyPayload, payloadInst)
			return next(c)
		}
	}
}
