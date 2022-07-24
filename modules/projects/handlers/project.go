package handlers

import (
	"fmt"
	"net/http"

	"github.com/dzungtran/echo-rest-api/modules/projects/domains"
	"github.com/dzungtran/echo-rest-api/modules/projects/dto"
	"github.com/dzungtran/echo-rest-api/modules/projects/usecases"
	"github.com/dzungtran/echo-rest-api/pkg/constants"
	"github.com/dzungtran/echo-rest-api/pkg/contexts"
	"github.com/dzungtran/echo-rest-api/pkg/logger"
	"github.com/dzungtran/echo-rest-api/pkg/middlewares"
	"github.com/dzungtran/echo-rest-api/pkg/utils"
	"github.com/dzungtran/echo-rest-api/pkg/wrapper"
	"github.com/labstack/echo/v4"
)

type ProjectHandler struct {
	ProjectUC usecases.ProjectUsecase
}

// NewProjectHandler will initialize the project resources endpoint
func NewProjectHandler(g *echo.Group, middManager *middlewares.MiddlewareManager, projectUsecase usecases.ProjectUsecase) {
	handler := &ProjectHandler{
		ProjectUC: projectUsecase,
	}

	apiV1 := g.Group("admin/projects", middManager.Auth())
	apiV1.GET("", wrapper.Wrap(handler.Fetch), middManager.CheckPoliciesWithRequestPayload(&dto.SearchProjectsReq{})).Name = "list:project"
	apiV1.POST("", wrapper.Wrap(handler.Create), middManager.CheckPoliciesWithRequestPayload(&dto.CreateProjectReq{})).Name = "create:project"

	apiV1Resource := g.Group("admin/projects/:projectId",
		middManager.Auth(),
		middlewares.RequireResourceIdInParam("projectId"),
		middManager.CheckPoliciesWithProject(),
	)
	apiV1Resource.GET("", wrapper.Wrap(handler.GetByID)).Name = "read:project"
	apiV1Resource.PUT("", wrapper.Wrap(handler.Update)).Name = "update:project"
	apiV1Resource.DELETE("", wrapper.Wrap(handler.Delete)).Name = "delete:project"
}

// Create will store the Project by given request body
func (h *ProjectHandler) Create(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()
	var err error

	payload := c.Get(constants.ContextKeyPayload)
	req, ok := payload.(*dto.CreateProjectReq)
	if !ok {
		return wrapper.Response{
			Status: http.StatusBadRequest,
			Error:  fmt.Errorf("invalid payload"),
		}
	}

	var proj *domains.Project
	if proj, err = h.ProjectUC.Create(ctx, *req); err != nil {
		logger.Log().Errorw("Error while create project", "error", err.Error())
		return wrapper.Response{
			Status: utils.GetHttpStatusCodeByErrorType(err, http.StatusInternalServerError),
			Error:  err,
		}
	}

	return wrapper.Response{Status: http.StatusCreated, Data: proj}
}

// GetByID will get Project by given id
func (h *ProjectHandler) GetByID(c echo.Context) wrapper.Response {
	project := contexts.GetProjectFromContext(c)
	return wrapper.Response{
		Data: project,
	}
}

// Fetch will fetch the Project
func (h *ProjectHandler) Fetch(c echo.Context) wrapper.Response {
	ctx := c.Request().Context()
	payload := c.Get(constants.ContextKeyPayload)
	req, ok := payload.(*dto.SearchProjectsReq)
	if !ok {
		return wrapper.Response{
			Status: http.StatusBadRequest,
			Error:  fmt.Errorf("invalid payload"),
		}
	}

	projects, count, err := h.ProjectUC.Fetch(ctx, *req)
	if err != nil {
		return wrapper.Response{
			Error:  err,
			Status: utils.GetHttpStatusCodeByErrorType(err, http.StatusInternalServerError),
		}
	}

	return wrapper.Response{
		Data:         projects,
		Total:        count,
		IncludeTotal: true,
	}
}

// Update will get project by given request body
func (h *ProjectHandler) Update(c echo.Context) wrapper.Response {
	var err error
	ctx := c.Request().Context()

	payload := c.Get(constants.ContextKeyPayload)
	req, ok := payload.(*dto.UpdateProjectReq)
	if !ok {
		return wrapper.Response{
			Status: http.StatusBadRequest,
			Error:  fmt.Errorf("invalid payload"),
		}
	}

	if err = h.ProjectUC.Update(ctx, req.ID, *req); err != nil {
		logger.Log().Errorw("Error while create project", "error", err.Error())
		if err == constants.ErrNotFound {
			return wrapper.Response{
				Status: http.StatusNotFound,
				Error:  utils.NewNotFoundError(),
			}
		}
		return wrapper.Response{
			Error:  err,
			Status: utils.GetHttpStatusCodeByErrorType(err, http.StatusInternalServerError),
		}
	}

	return wrapper.Response{}
}

// Delete will delete project by given param
func (h *ProjectHandler) Delete(c echo.Context) wrapper.Response {
	var err error
	ctx := c.Request().Context()
	id := utils.GetResourceIdFromParam(c, "projectId")

	if err = h.ProjectUC.Delete(ctx, id); err != nil {
		return wrapper.Response{
			Status: utils.GetHttpStatusCodeByErrorType(err, http.StatusInternalServerError),
			Error:  utils.NewError(err, ""),
		}
	}

	return wrapper.Response{}
}
