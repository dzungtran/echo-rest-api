package http

// Target: delivery/http/org_handler.go

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/dzungtran/echo-rest-api/delivery/requests"
	"github.com/dzungtran/echo-rest-api/delivery/wrapper"
	"github.com/dzungtran/echo-rest-api/domains"
	"github.com/dzungtran/echo-rest-api/pkg/authz"
	"github.com/dzungtran/echo-rest-api/pkg/logger"
	"github.com/dzungtran/echo-rest-api/pkg/middlewares"
	"github.com/dzungtran/echo-rest-api/pkg/utils"
	"github.com/dzungtran/echo-rest-api/usecases"
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

	apiV1 := g.Group("/orgs", middManager.Auth())
	apiV1.GET("", wrapper.Wrap(handler.Fetch)).Name = "list:org"
	apiV1.POST("", wrapper.Wrap(handler.Create)).Name = "create:org"

	apiV1.GET("/:orgId", wrapper.Wrap(handler.GetByID)).Name = "read:org"
	apiV1.PUT("/:orgId", wrapper.Wrap(handler.Update)).Name = "update:org"
	apiV1.DELETE("/:orgId", wrapper.Wrap(handler.Delete)).Name = "delete:org"

	apiV1.POST("/:orgId/invites", wrapper.Wrap(handler.Invite)).Name = "invite:org"
}

// Create will store the Org by given request body
func (h *OrgHandler) Create(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()
	var req requests.CreateOrgReq

	if err := c.Bind(&req); err != nil {
		return wrapper.Response{
			Status: http.StatusBadRequest,
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

	user, _ := authz.GetUserFromContext(c)
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
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("orgId"))
	if err != nil {
		return wrapper.Response{
			Status: http.StatusBadRequest,
			Error:  errors.New("invalid id"),
		}
	}

	denyMsg, err := authz.CheckPoliciesContext(c, authz.WithInputOrg(&domains.Org{
		Id: int64(id),
	}))
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

	org, err := h.OrgUC.GetByID(ctx, int64(id))
	if err != nil {
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

	return wrapper.Response{
		Data: org,
	}
}

// Fetch will fetch the Org
func (h *OrgHandler) Fetch(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()

	var req requests.SearchOrgsReq
	if err := c.Bind(&req); err != nil {
		return wrapper.Response{
			Status: http.StatusBadRequest,
			Error:  utils.NewError(err, ""),
		}
	}

	user, _ := authz.GetUserFromContext(c)
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

	var req requests.UpdateOrgReq
	if err = c.Bind(&req); err != nil {
		return wrapper.Response{
			Status: http.StatusUnprocessableEntity,
			Error:  utils.NewError(err, ""),
		}
	}

	denyMsg, err := authz.CheckPoliciesContext(c, authz.WithInputOrg(&domains.Org{
		Id: req.OrgId,
	}))
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

	denyMsg, err := authz.CheckPoliciesContext(c, authz.WithInputOrg(&domains.Org{
		Id: int64(id),
	}))
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
	var req requests.InviteUsers

	if err = c.Bind(&req); err != nil {
		return wrapper.Response{
			Status: http.StatusUnprocessableEntity,
			Error:  utils.NewError(err, ""),
		}
	}

	denyMsg, err := authz.CheckPoliciesContext(c, authz.WithInputOrg(&domains.Org{
		Id: req.OrgId,
	}))
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
