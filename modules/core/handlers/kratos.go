package handlers

import (
	"io/ioutil"
	"net/http"

	"github.com/dzungtran/echo-rest-api/modules/core/dto"
	"github.com/dzungtran/echo-rest-api/modules/core/usecases"
	"github.com/dzungtran/echo-rest-api/pkg/logger"
	"github.com/dzungtran/echo-rest-api/pkg/middlewares"
	"github.com/dzungtran/echo-rest-api/pkg/utils"
	"github.com/dzungtran/echo-rest-api/pkg/wrapper"
	"github.com/jinzhu/copier"
	"github.com/labstack/echo/v4"
	"github.com/tidwall/gjson"
)

type KratosHookHandler struct {
	UserUC usecases.UserUsecase
}

// NewKratosHookHandler will initialize the user resources endpoint
func NewKratosHookHandler(g *echo.Group, middManager *middlewares.MiddlewareManager, userUsecase usecases.UserUsecase) {
	handler := &KratosHookHandler{
		UserUC: userUsecase,
	}

	apiV1 := g.Group("hooks/kratos", middManager.KratosWebhookAuth)
	apiV1.POST("/after-registration", wrapper.Wrap(handler.AfterRegistration))
	apiV1.POST("/after-settings", wrapper.Wrap(handler.AfterSettings))
}

// Create will store the user by given request body
func (h *KratosHookHandler) AfterRegistration(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()
	resp, _ := ioutil.ReadAll(c.Request().Body)
	var err error
	req := dto.CreateUserReq{
		Email:     gjson.GetBytes(resp, "traits.email").String(),
		FirstName: gjson.GetBytes(resp, "traits.name.first").String(),
		LastName:  gjson.GetBytes(resp, "traits.name.last").String(),
		Code:      gjson.GetBytes(resp, "identity_id").String(),
	}

	if len(req.Email) == 0 {
		return wrapper.Response{
			Status: http.StatusBadRequest,
			Error:  utils.NewError(err, "invalid payload"),
		}
	}

	user, err := h.UserUC.GetByEmail(ctx, req.Email)
	if err == nil {
		// update user code
		updateReq := dto.UpdateUserReq{}
		copier.Copy(&updateReq, user)
		copier.Copy(&updateReq, req)

		err = h.UserUC.Update(ctx, user.Id, updateReq)
		if err != nil {
			logger.Log().Errorw("error while update user", "error", err)
			return wrapper.Response{
				Status: http.StatusInternalServerError,
				Error:  utils.NewError(err, "error while update user"),
			}
		}
		return wrapper.Response{Status: http.StatusOK}
	}

	_, err = h.UserUC.Register(ctx, req)
	if err != nil {
		if utils.IsCueError(err) {
			logger.Log().Debugw("invalid create user request", "error", err)
			return wrapper.Response{
				Status: http.StatusBadRequest,
				Error:  utils.NewError(err, "invalid payload"),
			}
		}

		logger.Log().Errorw("error while create user", "error", err)
		return wrapper.Response{
			Status: http.StatusInternalServerError,
			Error:  utils.NewError(err, "error while create user"),
		}
	}

	return wrapper.Response{Status: http.StatusCreated}
}

func (h *KratosHookHandler) AfterSettings(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()
	resp, _ := ioutil.ReadAll(c.Request().Body)
	var err error
	userId := int64(0)

	req := dto.UpdateUserReq{
		FirstName: gjson.GetBytes(resp, "traits.name.first").String(),
		LastName:  gjson.GetBytes(resp, "traits.name.last").String(),
	}

	if err = h.UserUC.Update(ctx, userId, req); err != nil {
		if utils.IsCueError(err) {
			logger.Log().Debugw("invalid update user request", "error", err)
			return wrapper.Response{
				Status: http.StatusBadRequest,
				Error:  utils.NewError(err, "invalid payload"),
			}
		}

		logger.Log().Errorw("error while create user", "error", err)
		return wrapper.Response{
			Status: http.StatusInternalServerError,
			Error:  utils.NewError(err, "error while update user"),
		}
	}

	return wrapper.Response{Status: http.StatusOK}
}
