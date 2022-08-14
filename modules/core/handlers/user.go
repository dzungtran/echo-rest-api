package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dzungtran/echo-rest-api/modules/core/domains"
	"github.com/dzungtran/echo-rest-api/modules/core/dto"
	"github.com/dzungtran/echo-rest-api/modules/core/usecases"
	"github.com/dzungtran/echo-rest-api/pkg/authz"
	"github.com/dzungtran/echo-rest-api/pkg/constants"
	"github.com/dzungtran/echo-rest-api/pkg/middlewares"
	"github.com/dzungtran/echo-rest-api/pkg/utils"
	"github.com/dzungtran/echo-rest-api/pkg/wrapper"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	UserUC usecases.UserUsecase
}

// NewUserHandler will initialize the user resources endpoint
func NewUserHandler(g *echo.Group, middManager *middlewares.MiddlewareManager, userUsecase usecases.UserUsecase) {
	handler := &UserHandler{
		UserUC: userUsecase,
	}

	apiV1 := g.Group("admin/users", middManager.Auth(), middManager.CheckPolicies())
	apiV1.GET("", wrapper.Wrap(handler.Fetch)).Name = "list:user"
	apiV1.POST("", wrapper.Wrap(handler.Create)).Name = "create:user"

	apiV1Resource := g.Group("admin/users/:userId", middManager.Auth())
	apiV1Resource.GET("", wrapper.Wrap(handler.GetByID)).Name = "read:user"
	apiV1Resource.PUT("", wrapper.Wrap(handler.Update), middManager.CheckPolicies()).Name = "update:user"
	apiV1Resource.DELETE("", wrapper.Wrap(handler.Delete), middManager.CheckPolicies()).Name = "delete:user"
}

// Create will store the user by given request body
func (h *UserHandler) Create(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()
	var req dto.CreateUserReq
	var user *domains.User
	var err error

	if err = c.Bind(&req); err != nil {
		return wrapper.Response{
			Status: http.StatusBadRequest,
			Error:  utils.NewError(err, ""),
		}
	}

	if user, err = h.UserUC.Create(ctx, req); err != nil {
		return wrapper.Response{
			Status: http.StatusInternalServerError,
			Error:  utils.NewError(err, ""),
		}
	}

	return wrapper.Response{Status: http.StatusCreated, Data: user}
}

// GetUserInfo godoc
// @Summary      Get user info
// @Description  Get user info ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        userId   path      int  true  "User ID"
// @Success      200  {object}  wrapper.SuccessResponse{data=domains.User}
// @Router       /users/{userId} [get]
func (h *UserHandler) GetByID(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		return wrapper.Response{
			Status: http.StatusBadRequest,
			Error:  errors.New("invalid id"),
		}
	}

	denyMsg, err := authz.CheckPoliciesContext(c, authz.WithInputExtraData("user_info", &domains.User{
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

	user, err := h.UserUC.GetByID(ctx, int64(id))
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
		Data: user,
	}
}

// Fetch will fetch the user
func (h *UserHandler) Fetch(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()

	var req dto.SearchUsersReq
	if err := c.Bind(&req); err != nil {
		return wrapper.Response{
			Status: http.StatusBadRequest,
			Error:  utils.NewError(err, ""),
		}
	}

	users, count, err := h.UserUC.Fetch(ctx, req)
	if err != nil {
		return wrapper.Response{
			Error:  err,
			Status: http.StatusInternalServerError,
		}
	}

	return wrapper.Response{
		Data:         users,
		Total:        count,
		IncludeTotal: true,
	}
}

// Update will get user by given request body
func (h *UserHandler) Update(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		return wrapper.Response{
			Status: http.StatusBadRequest,
			Error:  utils.NewNotFoundError(),
		}
	}

	var req dto.UpdateUserReq
	if err = c.Bind(&req); err != nil {
		return wrapper.Response{
			Status: http.StatusUnprocessableEntity,
			Error:  utils.NewError(err, ""),
		}
	}

	denyMsg, err := authz.CheckPoliciesContext(c, authz.WithInputExtraData("user_info", &domains.User{
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

	if err := h.UserUC.Update(ctx, int64(id), req); err != nil {
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

// Delete will delete user by given param
func (h *UserHandler) Delete(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		return wrapper.Response{
			Status: http.StatusBadRequest,
			Error:  utils.NewNotFoundError(),
		}
	}

	if err = h.UserUC.Delete(ctx, int64(id)); err != nil {
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
