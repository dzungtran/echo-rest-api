package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/dzungtran/echo-rest-api/modules/core/domains"
	"github.com/dzungtran/echo-rest-api/modules/core/dto"
	"github.com/dzungtran/echo-rest-api/pkg/authz"
	"github.com/dzungtran/echo-rest-api/pkg/constants"
	"github.com/dzungtran/echo-rest-api/pkg/utils"
	"github.com/dzungtran/echo-rest-api/pkg/wrapper"
	"github.com/labstack/echo/v4"
)

// GetUserInfo godoc
// @Summary      Get user info
// @Description  Get user info by ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        userId   path      int  true  "User ID"
// @Success      200  {object}  wrapper.SuccessResponse{data=domains.User}
// @Failure      400  {object}  wrapper.FailResponse
// @Failure      401  {object}  wrapper.FailResponse
// @Failure      403  {object}  wrapper.FailResponse
// @Failure      404  {object}  wrapper.FailResponse
// @Security     XFirebaseBearer
// @Router       /admin/users/{userId} [get]
func (h *UserHandler) GetByID(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()
	id, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		return wrapper.Response{
			Status: http.StatusBadRequest,
			Error:  utils.NewError(errors.New("invalid id"), ""),
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

// GetListUser godoc
// @Summary      Get list user info
// @Description  Get list user info
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        limit   query     int  false  "Number of records should be returned"
// @Param        page    query     int  false  "Page"
// @Success      200  {object}  wrapper.SuccessResponse{data=[]domains.User}
// @Failure      400  {object}  wrapper.FailResponse
// @Failure      401  {object}  wrapper.FailResponse
// @Failure      403  {object}  wrapper.FailResponse
// @Security     XFirebaseBearer
// @Router       /admin/users [get]
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
