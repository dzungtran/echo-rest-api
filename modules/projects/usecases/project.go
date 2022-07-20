package usecases

import (
	"context"

	"github.com/dzungtran/echo-rest-api/modules/projects/domains"
	"github.com/dzungtran/echo-rest-api/modules/projects/dto"
	"github.com/dzungtran/echo-rest-api/modules/projects/pkg/cue"
	"github.com/dzungtran/echo-rest-api/modules/projects/repositories"
	"github.com/dzungtran/echo-rest-api/pkg/contexts"
	"github.com/dzungtran/echo-rest-api/pkg/utils"
	"github.com/jinzhu/copier"
)

// ProjectUsecase represent the project's usecase contract
type ProjectUsecase interface {
	Create(ctx context.Context, request dto.CreateProjectReq) (*domains.Project, error)
	GetByID(ctx context.Context, id int64) (*domains.Project, error)
	Fetch(ctx context.Context, req dto.SearchProjectsReq) ([]*domains.Project, int64, error)
	Update(ctx context.Context, id int64, request dto.UpdateProjectReq) error
	Delete(ctx context.Context, id int64) error
}

type projectUsecase struct {
	projectRepo repositories.ProjectRepository
}

// NewProjectUsecase will create new an projectUsecase object representation of ProjectUsecase interface
func NewProjectUsecase(projectRepo repositories.ProjectRepository) ProjectUsecase {
	return &projectUsecase{
		projectRepo: projectRepo,
	}
}

func (u *projectUsecase) Create(ctx context.Context, req dto.CreateProjectReq) (project *domains.Project, err error) {
	project = &domains.Project{}
	copier.Copy(project, req)

	if err = utils.CueValidateObject("CreateProjectReq", cue.CueDefinitionForProject, req); err != nil {
		return nil, err
	}

	if project.Code == "" {
		project.Code = utils.GenerateLongUUID()
	}

	projectId, err := u.projectRepo.Create(ctx, project)
	if err != nil {
		return
	}

	project, err = u.projectRepo.GetByID(ctx, projectId)
	return
}

func (u *projectUsecase) GetByID(ctx context.Context, id int64) (project *domains.Project, err error) {
	project, err = u.projectRepo.GetByID(ctx, id)
	return
}

func (u *projectUsecase) Fetch(ctx context.Context, req dto.SearchProjectsReq) (projects []*domains.Project, count int64, err error) {
	p := repositories.ParamsForFetchProjects{
		CommonParamsForFetch: contexts.CommonParamsForFetch{
			Page:  uint64(req.Page),
			Limit: uint64(req.Limit),
		},
	}

	projects, count, err = u.projectRepo.Fetch(ctx, p)
	if err != nil {
		return
	}

	return
}

func (u *projectUsecase) Update(ctx context.Context, id int64, req dto.UpdateProjectReq) (err error) {
	project, err := u.projectRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	copier.Copy(project, req)
	err = u.projectRepo.Update(ctx, project, []string{"name", "description", "timezone", "default_language"})
	return
}

func (u *projectUsecase) Delete(ctx context.Context, id int64) (err error) {
	_, err = u.projectRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	err = u.projectRepo.DeleteById(ctx, id)
	return
}
