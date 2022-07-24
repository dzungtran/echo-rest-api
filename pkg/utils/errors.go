package utils

import (
	"net/http"

	cueErrors "cuelang.org/go/cue/errors"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

type Error struct {
	Errors map[string]interface{} `json:"errors"`
}

func (e Error) Error() string {
	return e.Errors["message"].(string)
}

func (e Error) RawError() error {
	return e.Errors["raw"].(error)
}

// NewError creates a new error response
func NewError(err error, msg string) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["raw"] = err

	switch v := err.(type) {
	case *echo.HTTPError:
		e.Errors["message"] = v.Message
	default:
		if msg != "" {
			e.Errors["message"] = msg
		} else {
			e.Errors["message"] = v.Error()
		}
	}
	return e
}

type invalidField struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// NewValidationError creates a new error response representing a data validation error (HTTP 400)
func NewValidationError(errs validator.ValidationErrors) Error {
	e := Error{}
	e.Errors = make(map[string]interface{})

	var details []invalidField
	for _, field := range errs {
		details = append(details, invalidField{
			Field: field.Field(),
			Error: field.Error(),
		})
	}

	e.Errors["message"] = "there is some problem with the data you submitted"
	e.Errors["details"] = details

	return e
}

// NewAccessForbiddenError creates a new error response representing an authorization failure (HTTP 403)
func NewAccessForbiddenError() Error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["message"] = "access forbidden"
	return e
}

// NewNotFoundError creates a new error response representing a resource-not-found error (HTTP 404)
func NewNotFoundError() error {
	e := Error{}
	e.Errors = make(map[string]interface{})
	e.Errors["message"] = "resource not found"
	return e
}

func IsCueError(err error) bool {
	_, ok := err.(cueErrors.Error)
	return ok
}

func IsDuplicatedError(err error) bool {
	switch err.(type) {
	case *pq.Error:
		return err.(pq.Error).Code == "23505"
	}

	return false
}

func GetHttpStatusCodeByErrorType(err error, defaultCode int) int {
	switch err.(type) {
	case *pq.Error:
		return http.StatusInternalServerError
	case cueErrors.Error:
	case validator.ValidationErrors:
		return http.StatusBadRequest

	}

	if defaultCode > 0 {
		return defaultCode
	}

	return http.StatusInternalServerError
}
