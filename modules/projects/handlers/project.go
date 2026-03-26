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
	// Check comment of function CheckPoliciesWithRequestPayload
	apiV1.GET("", wrapper.Wrap(handler.Fetch), middManager.CheckPoliciesWithRequestPayload(new(dto.SearchProjectsReq))).Name = "list:project"
	apiV1.POST("", wrapper.Wrap(handler.Create), middManager.CheckPoliciesWithRequestPayload(new(dto.CreateProjectReq))).Name = "create:project"

	apiV1Resource := g.Group("admin/projects/:projectId",
		middManager.Auth(),
		middlewares.RequireResourceIdInParam("projectId"),
		middManager.CheckPoliciesWithProject(),
	)
	apiV1Resource.GET("", wrapper.Wrap(handler.GetByID)).Name = "read:project"
	apiV1Resource.PUT("", wrapper.Wrap(handler.Update)).Name = "update:project"
	apiV1Resource.DELETE("", wrapper.Wrap(handler.Delete)).Name = "delete:project"
}

// CreateProject godoc
// @Summary      Create a new project
// @Description  Create a new project under an organization
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        body    body      dto.CreateProjectReq  true  "Project creation request"
// @Success      201  {object}  wrapper.SuccessResponse{data=domains.Project}
// @Failure      400  {object}  wrapper.FailResponse
// @Failure      401  {object}  wrapper.FailResponse
// @Failure      403  {object}  wrapper.FailResponse
// @Failure      500  {object}  wrapper.FailResponse
// @Security     XFirebaseBearer
// @Router       /admin/projects [post]
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

// GetProjectInfo godoc
// @Summary      Get project info
// @Description  Get project by ID
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        projectId   path      int  true  "Project ID"
// @Success      200  {object}  wrapper.SuccessResponse{data=domains.Project}
// @Failure      401  {object}  wrapper.FailResponse
// @Failure      403  {object}  wrapper.FailResponse
// @Failure      404  {object}  wrapper.FailResponse
// @Security     XFirebaseBearer
// @Router       /admin/projects/{projectId} [get]
func (h *ProjectHandler) GetByID(c echo.Context) wrapper.Response {
	project := contexts.GetProjectFromContext(c)
	return wrapper.Response{
		Data: project,
	}
}

// ListProjects godoc
// @Summary      List projects
// @Description  Get list of projects
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        limit   query     int  false  "Number of records should be returned"
// @Param        page    query     int  false  "Page"
// @Success      200  {object}  wrapper.SuccessResponse{data=[]domains.Project}
// @Failure      400  {object}  wrapper.FailResponse
// @Failure      401  {object}  wrapper.FailResponse
// @Failure      403  {object}  wrapper.FailResponse
// @Security     XFirebaseBearer
// @Router       /admin/projects [get]
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

// UpdateProject godoc
// @Summary      Update project
// @Description  Update project by ID
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        projectId   path      int  true  "Project ID"
// @Param        body    body      dto.UpdateProjectReq  true  "Project update request"
// @Success      200  {object}  wrapper.SuccessResponse{}
// @Failure      400  {object}  wrapper.FailResponse
// @Failure      401  {object}  wrapper.FailResponse
// @Failure      403  {object}  wrapper.FailResponse
// @Failure      404  {object}  wrapper.FailResponse
// @Failure      500  {object}  wrapper.FailResponse
// @Security     XFirebaseBearer
// @Router       /admin/projects/{projectId} [put]
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

// DeleteProject godoc
// @Summary      Delete project
// @Description  Delete project by ID
// @Tags         projects
// @Accept       json
// @Produce      json
// @Param        projectId   path      int  true  "Project ID"
// @Success      200  {object}  wrapper.SuccessResponse{}
// @Failure      400  {object}  wrapper.FailResponse
// @Failure      401  {object}  wrapper.FailResponse
// @Failure      403  {object}  wrapper.FailResponse
// @Failure      404  {object}  wrapper.FailResponse
// @Failure      500  {object}  wrapper.FailResponse
// @Security     XFirebaseBearer
// @Router       /admin/projects/{projectId} [delete]
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
