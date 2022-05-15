package http

// Target: delivery/http/{{ .ModuleName | ToSnake }}_handler.go

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"{{ .RootPackage }}/delivery/requests"
	"{{ .RootPackage }}/delivery/wrapper"
	"{{ .RootPackage }}/domains"
	"{{ .RootPackage }}/pkg/authz"
	"{{ .RootPackage }}/pkg/constants"
	"{{ .RootPackage }}/pkg/middlewares"
	"{{ .RootPackage }}/pkg/utils"
	"{{ .RootPackage }}/usecases"
)

type {{ .ModuleName }}Handler struct {
	{{ .ModuleName }}UC usecases.{{ .ModuleName }}Usecase
}

// New{{ .ModuleName }}Handler will initialize the {{ .ModuleName | ToLower }} resources endpoint
func New{{ .ModuleName }}Handler(g *echo.Group, middManager *middlewares.MiddlewareManager, {{ .ModuleName | ToLowerCamel }}Usecase usecases.{{ .ModuleName }}Usecase) {
	handler := &{{ .ModuleName }}Handler{
		{{ .ModuleName }}UC: {{ .ModuleName | ToLowerCamel }}Usecase,
	}

	apiV1 := g.Group("/{{ .ModuleName | ToKebab }}s")
	apiV1.GET("", wrapper.Wrap(handler.Fetch)).Name = "list:{{ .ModuleName | ToLowerCamel }}"
	apiV1.POST("", wrapper.Wrap(handler.Create)).Name = "create:{{ .ModuleName | ToLowerCamel }}"
	
	apiV1.GET("/:{{ .ModuleName | ToLowerCamel }}Id", wrapper.Wrap(handler.GetByID)).Name = "read:{{ .ModuleName | ToLowerCamel }}"
	apiV1.PUT("/:{{ .ModuleName | ToLowerCamel }}Id", wrapper.Wrap(handler.Update)).Name = "update:{{ .ModuleName | ToLowerCamel }}"
	apiV1.DELETE("/:{{ .ModuleName | ToLowerCamel }}Id", wrapper.Wrap(handler.Delete)).Name = "delete:{{ .ModuleName | ToLowerCamel }}"
}

// Create will store the {{ .ModuleName }} by given request body
func (h *{{ .ModuleName }}Handler) Create(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()
	var req requests.Create{{ .ModuleName }}Req

	denyMsg, err := authz.CheckPoliciesContext(c)
	if err != nil {
		msg := ""
		if len(denyMsg) > 0 {
			msg = denyMsg[0]
		}
		return wrapper.Response{
			Status: http.StatusForbidden,
			Error:  utils.NewError(err, msg),
		}
	}

	if err := c.Bind(&req); err != nil {
		return wrapper.Response{
			Status: http.StatusBadRequest,
			Error:  utils.NewError(err, ""),
		}
	}

	var new{{ .ModuleName }} *domains.{{ .ModuleName }}
	if new{{ .ModuleName }}, err = h.{{ .ModuleName }}UC.Create(ctx, req); err != nil {
		return wrapper.Response{
			Status: http.StatusInternalServerError,
			Error:  utils.NewError(err, ""),
		}
	}

	return wrapper.Response{Status: http.StatusCreated, Data: new{{ .ModuleName }}}
}

// GetByID will get {{ .ModuleName }} by given id
func (h *{{ .ModuleName }}Handler) GetByID(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("{{ .ModuleName | ToLowerCamel }}Id"))
	if err != nil {
		return wrapper.Response{
			Status: http.StatusBadRequest,
			Error:  errors.New("invalid id"),
		}
	}

	denyMsg, err := authz.CheckPoliciesContext(c)
	if err != nil {
		msg := ""
		if len(denyMsg) > 0 {
			msg = denyMsg[0]
		}
		return wrapper.Response{
			Status: http.StatusForbidden,
			Error:  utils.NewError(err, msg),
		}
	}

	{{ .ModuleName | ToLowerCamel }}, err := h.{{ .ModuleName }}UC.GetByID(ctx, int64(id))
	if err != nil {
		if err == constants.ErrNotFound {
			return wrapper.Response{
				Status: http.StatusNotFound,
				Error:  utils.NewNotFoundError(),
			}
		}
		return wrapper.Response{
			Status: http.StatusInternalServerError,
			Error:  utils.NewError(err, ""),
		}
	}

	return wrapper.Response{
		Data: {{ .ModuleName | ToLowerCamel }},
	}
}

// Fetch will fetch the {{ .ModuleName }}
func (h *{{ .ModuleName }}Handler) Fetch(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()

	denyMsg, err := authz.CheckPoliciesContext(c)
	if err != nil {
		msg := ""
		if len(denyMsg) > 0 {
			msg = denyMsg[0]
		}
		return wrapper.Response{
			Status: http.StatusForbidden,
			Error:  utils.NewError(err, msg),
		}
	}

	var req requests.Search{{ .ModuleName }}sReq
	if err := c.Bind(&req); err != nil {
		return wrapper.Response{
			Status: http.StatusBadRequest,
			Error:  utils.NewError(err, ""),
		}
	}

	if err := c.Validate(req); err != nil {
		errValidations := err.(validator.ValidationErrors)
		return wrapper.Response{
			Status: http.StatusBadRequest,
			Error:  utils.NewValidationError(errValidations),
		}
	}

	{{ .ModuleName | ToLower }}s, count, err := h.{{ .ModuleName }}UC.Fetch(ctx, req)
	if err != nil {
		return wrapper.Response{
			Error:  err,
			Status: http.StatusInternalServerError,
		}
	}

	return wrapper.Response{
		Data:         {{ .ModuleName | ToLower }}s,
		Total:        count,
		IncludeTotal: true,
	}
}

// Update will get {{ .ModuleName | ToLower }} by given request body
func (h *{{ .ModuleName }}Handler) Update(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()
	var err error

	var req requests.Update{{ .ModuleName }}Req
	if err = c.Bind(&req); err != nil {
		return wrapper.Response{
			Status: http.StatusUnprocessableEntity,
			Error:  utils.NewError(err, ""),
		}
	}

	denyMsg, err := authz.CheckPoliciesContext(c)
	if err != nil {
		msg := ""
		if len(denyMsg) > 0 {
			msg = denyMsg[0]
		}
		return wrapper.Response{
			Status: http.StatusForbidden,
			Error:  utils.NewError(err, msg),
		}
	}

	if err := h.{{ .ModuleName }}UC.Update(ctx, req.ID, req); err != nil {
		if err == constants.ErrNotFound {
			return wrapper.Response{
				Status: http.StatusNotFound,
				Error:  utils.NewNotFoundError(),
			}
		}
		return wrapper.Response{
			Error:  err,
			Status: http.StatusInternalServerError,
		}
	}

	return wrapper.Response{}
}

// Delete will delete {{ .ModuleName | ToLower }} by given param
func (h *{{ .ModuleName }}Handler) Delete(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("{{ .ModuleName | ToLowerCamel }}Id"))
	if err != nil {
		return wrapper.Response{
			Status: http.StatusBadRequest,
			Error:  utils.NewNotFoundError(),
		}
	}

	denyMsg, err := authz.CheckPoliciesContext(c)
	if err != nil {
		msg := ""
		if len(denyMsg) > 0 {
			msg = denyMsg[0]
		}
		return wrapper.Response{
			Status: http.StatusForbidden,
			Error:  utils.NewError(err, msg),
		}
	}

	if err = h.{{ .ModuleName }}UC.Delete(ctx, int64(id)); err != nil {
		if err == constants.ErrNotFound {
			return wrapper.Response{
				Status: http.StatusNotFound,
				Error:  utils.NewNotFoundError(),
			}
		}
		return wrapper.Response{
			Status: http.StatusInternalServerError,
			Error:  utils.NewError(err, ""),
		}
	}

	return wrapper.Response{}
}