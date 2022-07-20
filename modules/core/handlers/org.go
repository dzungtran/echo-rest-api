package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/dzungtran/echo-rest-api/modules/core/domains"
	"github.com/dzungtran/echo-rest-api/modules/core/dto"
	"github.com/dzungtran/echo-rest-api/modules/core/usecases"
	"github.com/dzungtran/echo-rest-api/pkg/contexts"
	"github.com/dzungtran/echo-rest-api/pkg/logger"
	"github.com/dzungtran/echo-rest-api/pkg/middlewares"
	"github.com/dzungtran/echo-rest-api/pkg/utils"
	"github.com/dzungtran/echo-rest-api/pkg/wrapper"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type OrgHandler struct {
	OrgUC usecases.OrgUsecase
}

// NewOrgHandler will initialize the org resources endpoint
func NewOrgHandler(g *echo.Group, middManager *middlewares.MiddlewareManager, orgUsecase usecases.OrgUsecase) {
	handler := &OrgHandler{
		OrgUC: orgUsecase,
	}

	apiV1 := g.Group("admin/orgs", middManager.Auth(), middManager.CheckPolicies())
	apiV1.GET("", wrapper.Wrap(handler.Fetch)).Name = "list:org"
	apiV1.POST("", wrapper.Wrap(handler.Create)).Name = "create:org"

	apiV1Resource := g.Group("admin/orgs/:orgId",
		middManager.Auth(),
		middlewares.RequireResourceIdInParam("orgId"),
		middManager.CheckPoliciesWithOrg(),
	)
	apiV1Resource.GET("", wrapper.Wrap(handler.GetByID)).Name = "read:org"
	apiV1Resource.PUT("", wrapper.Wrap(handler.Update)).Name = "update:org"
	apiV1Resource.DELETE("", wrapper.Wrap(handler.Delete)).Name = "delete:org"

	apiV1Resource.POST("/invites", wrapper.Wrap(handler.Invite)).Name = "invite:org"
}

// Create will store the Org by given request body
func (h *OrgHandler) Create(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()
	var req dto.CreateOrgReq
	var err error

	if err := c.Bind(&req); err != nil {
		return wrapper.Response{
			Status: http.StatusBadRequest,
			Error:  utils.NewError(err, ""),
		}
	}

	user, _ := contexts.GetUserFromContext(c)
	req.UserId = user.Id

	var newOrg *domains.Org
	if newOrg, err = h.OrgUC.Create(ctx, req); err != nil {
		if utils.IsCueError(err) {
			logger.Log().Debugw("invalid create org request", "error", err)
			return wrapper.Response{
				Status: http.StatusBadRequest,
				Error:  utils.NewError(err, "invalid payload"),
			}
		}

		logger.Log().Errorw("error while create org", "error", err)
		return wrapper.Response{
			Status: http.StatusInternalServerError,
			Error:  utils.NewError(err, "internal server error"),
		}
	}

	return wrapper.Response{Status: http.StatusCreated, Data: newOrg}
}

// GetByID will get Org by given id
func (h *OrgHandler) GetByID(c echo.Context) wrapper.Response {
	org := contexts.GetOrgFromContext(c)
	return wrapper.Response{
		Data: org,
	}
}

// Fetch will fetch the Org
func (h *OrgHandler) Fetch(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()

	var req dto.SearchOrgsReq
	if err := c.Bind(&req); err != nil {
		return wrapper.Response{
			Status: http.StatusBadRequest,
			Error:  utils.NewError(err, ""),
		}
	}

	user, _ := contexts.GetUserFromContext(c)
	if err := c.Validate(req); err != nil {
		errValidations := err.(validator.ValidationErrors)
		return wrapper.Response{
			Status: http.StatusBadRequest,
			Error:  utils.NewValidationError(errValidations),
		}
	}

	// Just allows get org which user opt in
	req.Ids = user.GetOrgIds()
	if len(req.Ids) == 0 {
		return wrapper.Response{
			Data:         []interface{}{},
			Total:        0,
			IncludeTotal: true,
		}
	}

	orgs, count, err := h.OrgUC.Fetch(ctx, req)
	if err != nil {
		return wrapper.Response{
			Error:  err,
			Status: http.StatusInternalServerError,
		}
	}

	return wrapper.Response{
		Data:         orgs,
		Total:        count,
		IncludeTotal: true,
	}
}

// Update will get org by given request body
func (h *OrgHandler) Update(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()
	var err error

	var req dto.UpdateOrgReq
	if err = c.Bind(&req); err != nil {
		return wrapper.Response{
			Status: http.StatusUnprocessableEntity,
			Error:  utils.NewError(err, ""),
		}
	}

	if err := h.OrgUC.Update(ctx, req.OrgId, req); err != nil {
		if utils.IsCueError(err) {
			logger.Log().Debugw("invalid update org request", "error", err)
			return wrapper.Response{
				Status: http.StatusBadRequest,
				Error:  utils.NewError(err, "invalid payload"),
			}
		}

		if err == sql.ErrNoRows {
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

// Delete will delete org by given param
func (h *OrgHandler) Delete(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("orgId"))
	if err != nil {
		return wrapper.Response{
			Status: http.StatusBadRequest,
			Error:  utils.NewNotFoundError(),
		}
	}

	if err = h.OrgUC.Delete(ctx, int64(id)); err != nil {
		if err == sql.ErrNoRows {
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

func (h *OrgHandler) Invite(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()
	var err error
	var req dto.InviteUsers

	if err = c.Bind(&req); err != nil {
		return wrapper.Response{
			Status: http.StatusUnprocessableEntity,
			Error:  utils.NewError(err, ""),
		}
	}

	if err := h.OrgUC.Invite(ctx, req.OrgId, req); err != nil {
		if utils.IsCueError(err) {
			logger.Log().Debugw("invalid invite to org request", "error", err)
			return wrapper.Response{
				Status: http.StatusBadRequest,
				Error:  utils.NewError(err, "invalid payload"),
			}
		}

		return wrapper.Response{
			Error:  err,
			Status: http.StatusInternalServerError,
		}
	}

	return wrapper.Response{}
}
