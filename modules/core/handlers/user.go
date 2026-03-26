package handlers

import (
	"github.com/dzungtran/echo-rest-api/modules/core/usecases"
	"github.com/dzungtran/echo-rest-api/pkg/contexts"
	"github.com/dzungtran/echo-rest-api/pkg/middlewares"
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

	apiMeV1 := g.Group("me", middManager.Auth())
	apiMeV1.GET("", wrapper.Wrap(handler.GetCurrentUserInfo)).Name = "read:me"
	apiMeV1.PUT("", wrapper.Wrap(handler.UpdateCurrentUserInfo), middManager.CheckPolicies()).Name = "update:me"

	// Endpoints for user management
	apiV1 := g.Group("admin/users", middManager.Auth(), middManager.CheckPolicies())
	apiV1.GET("", wrapper.Wrap(handler.Fetch)).Name = "list:user"
	// apiV1.POST("", wrapper.Wrap(handler.Create)).Name = "create:user"

	apiV1Resource := g.Group("admin/users/:userId", middManager.Auth())
	apiV1Resource.GET("", wrapper.Wrap(handler.GetByID)).Name = "read:user"
	// apiV1Resource.PUT("", wrapper.Wrap(handler.Update), middManager.CheckPolicies()).Name = "update:user"
	// apiV1Resource.DELETE("", wrapper.Wrap(handler.Delete), middManager.CheckPolicies()).Name = "delete:user"
}

// GetCurrentUserInfo godoc
// @Summary      Get current user info
// @Description  Get current authenticated user info
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  wrapper.SuccessResponse{data=domains.User}
// @Failure      401  {object}  wrapper.FailResponse
// @Security     XFirebaseBearer
// @Router       /me [get]
func (h *UserHandler) GetCurrentUserInfo(c echo.Context) wrapper.Response {
	user, _ := contexts.GetUserFromContext(c)
	return wrapper.Response{
		Data: user,
	}
}

// UpdateCurrentUserInfo godoc
// @Summary      Update current user info
// @Description  Update current authenticated user info
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  wrapper.SuccessResponse{data=domains.User}
// @Failure      400  {object}  wrapper.FailResponse
// @Failure      401  {object}  wrapper.FailResponse
// @Failure      403  {object}  wrapper.FailResponse
// @Security     XFirebaseBearer
// @Router       /me [put]
func (h *UserHandler) UpdateCurrentUserInfo(c echo.Context) wrapper.Response {
	user, _ := contexts.GetUserFromContext(c)

	// TODO: update user info here

	return wrapper.Response{
		Data: user,
	}
}
