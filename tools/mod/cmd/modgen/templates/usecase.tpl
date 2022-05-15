package usecases

// Target: usecases/{{ .ModuleName | ToSnake }}.go

import (
	"context"

	"github.com/jinzhu/copier"
	"{{ .RootPackage }}/delivery/defines"
	"{{ .RootPackage }}/delivery/requests"
	"{{ .RootPackage }}/domains"
	"{{ .RootPackage }}/repositories/postgres"
)

// {{ .ModuleName }}Usecase represent the {{ .ModuleName | ToLowerCamel }}'s usecase contract
type {{ .ModuleName }}Usecase interface {
	Create(ctx context.Context, request requests.Create{{ .ModuleName }}Req) (*domains.{{ .ModuleName }}, error)
	GetByID(ctx context.Context, id int64) (*domains.{{ .ModuleName }}, error)
	Fetch(ctx context.Context, req requests.Search{{ .ModuleName }}sReq) ([]*domains.{{ .ModuleName }}, int64, error)
	Update(ctx context.Context, id int64, request requests.Update{{ .ModuleName }}Req) error
	Delete(ctx context.Context, id int64) error
}

type {{ .ModuleName | ToLowerCamel }}Usecase struct {
	{{ .ModuleName | ToLowerCamel }}Repo postgres.{{ .ModuleName }}Repository
}

// New{{ .ModuleName }}Usecase will create new an {{ .ModuleName | ToLowerCamel }}Usecase object representation of {{ .ModuleName }}Usecase interface
func New{{ .ModuleName }}Usecase({{ .ModuleName | ToLowerCamel }}Repo postgres.{{ .ModuleName }}Repository) {{ .ModuleName }}Usecase {
	return &{{ .ModuleName | ToLowerCamel }}Usecase{
		{{ .ModuleName | ToLowerCamel }}Repo: {{ .ModuleName | ToLowerCamel }}Repo,
	}
}

func (u *{{ .ModuleName | ToLowerCamel }}Usecase) Create(ctx context.Context, req requests.Create{{ .ModuleName }}Req) ({{ .ModuleName | ToLowerCamel }} *domains.{{ .ModuleName }}, err error) {
	{{ .ModuleName | ToLowerCamel }} = &domains.{{ .ModuleName }}{}
	copier.Copy({{ .ModuleName | ToLowerCamel }}, req)

	//if err = utils.CueValidateObject("Create{{ .ModuleName }}Req", cue.CueDefinitionFor{{ .ModuleName }}, req); err != nil {
	//	return nil, err
	//}

	{{ .ModuleName | ToLowerCamel }}Id, err := u.{{ .ModuleName | ToLowerCamel }}Repo.Create(ctx, {{ .ModuleName | ToLowerCamel }})
	if err != nil {
		return
	}

	{{ .ModuleName | ToLowerCamel }}, err = u.{{ .ModuleName | ToLowerCamel }}Repo.GetByID(ctx, {{ .ModuleName | ToLowerCamel }}Id)
	return
}

func (u *{{ .ModuleName | ToLowerCamel }}Usecase) GetByID(ctx context.Context, id int64) ({{ .ModuleName | ToLowerCamel }} *domains.{{ .ModuleName }}, err error) {
	{{ .ModuleName | ToLowerCamel }}, err = u.{{ .ModuleName | ToLowerCamel }}Repo.GetByID(ctx, id)
	return
}

func (u *{{ .ModuleName | ToLowerCamel }}Usecase) Fetch(ctx context.Context, req requests.Search{{ .ModuleName }}sReq) ({{ .ModuleName | ToLowerCamel }}s []*domains.{{ .ModuleName }}, count int64, err error) {
	p := postgres.ParamsForFetch{{ .ModuleName }}s{
		CommonParamsForFetch: defines.CommonParamsForFetch{
			Page:  uint64(req.Page),
			Limit: uint64(req.Limit),
		},
	}

	{{ .ModuleName | ToLowerCamel }}s, count, err = u.{{ .ModuleName | ToLowerCamel }}Repo.Fetch(ctx, p)
	if err != nil {
		return
	}

	return
}

func (u *{{ .ModuleName | ToLowerCamel }}Usecase) Update(ctx context.Context, id int64, req requests.Update{{ .ModuleName }}Req) (err error) {
	{{ .ModuleName | ToLowerCamel }}, err := u.{{ .ModuleName | ToLowerCamel }}Repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	copier.Copy({{ .ModuleName | ToLowerCamel }}, req)
	err = u.{{ .ModuleName | ToLowerCamel }}Repo.Update(ctx, {{ .ModuleName | ToLowerCamel }}, []string{"first_name", "last_name", "email", "phone", "status"})
	return
}

func (u *{{ .ModuleName | ToLowerCamel }}Usecase) Delete(ctx context.Context, id int64) (err error) {
	_, err = u.{{ .ModuleName | ToLowerCamel }}Repo.GetByID(ctx, id)
	if err != nil {
		return
	}

	err = u.{{ .ModuleName | ToLowerCamel }}Repo.DeleteById(ctx, id)
	return
}
