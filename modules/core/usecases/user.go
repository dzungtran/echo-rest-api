package usecases

import (
	"context"

	"github.com/dzungtran/echo-rest-api/modules/core/domains"
	"github.com/dzungtran/echo-rest-api/modules/core/dto"
	"github.com/dzungtran/echo-rest-api/modules/core/repositories"
	"github.com/dzungtran/echo-rest-api/pkg/contexts"
	"github.com/dzungtran/echo-rest-api/pkg/cue"
	"github.com/dzungtran/echo-rest-api/pkg/utils"
	"github.com/jinzhu/copier"
)

// UserUsecase represent the user's usecase contract
type UserUsecase interface {
	Create(ctx context.Context, request dto.CreateUserReq) (*domains.User, error)
	GetByID(ctx context.Context, id int64) (*domains.User, error)
	Fetch(ctx context.Context, req dto.SearchUsersReq) ([]*domains.User, int64, error)
	Update(ctx context.Context, id int64, request dto.UpdateUserReq) error
	Delete(ctx context.Context, id int64) error
	GetByCode(ctx context.Context, code string) (*domains.User, error)
	GetByEmail(ctx context.Context, email string) (*domains.User, error)
	Register(ctx context.Context, request dto.CreateUserReq) (*domains.User, error)
}

type userUsecase struct {
	userRepo   repositories.UserRepository
	orgUsecase OrgUsecase
}

// NewUserUsecase will create new an userUsecase object representation of UserUsecase interface
func NewUserUsecase(userRepo repositories.UserRepository, orgUsecase OrgUsecase) UserUsecase {
	return &userUsecase{
		userRepo:   userRepo,
		orgUsecase: orgUsecase,
	}
}

func (u *userUsecase) Create(ctx context.Context, req dto.CreateUserReq) (user *domains.User, err error) {
	user = &domains.User{}
	copier.Copy(user, req)

	if user.Status == "" {
		user.Status = domains.UserStatusActive
	}

	if err = utils.CueValidateObject("CreateUserRequest", cue.CueDefinitionForUser, req); err != nil {
		return nil, err
	}

	userId, err := u.userRepo.Create(ctx, user)
	if err != nil {
		return
	}

	user, err = u.userRepo.GetByID(ctx, userId)
	return
}

func (u *userUsecase) GetByID(ctx context.Context, id int64) (user *domains.User, err error) {
	user, err = u.userRepo.GetByID(ctx, id)
	return
}

func (u *userUsecase) Fetch(ctx context.Context, req dto.SearchUsersReq) (users []*domains.User, count int64, err error) {
	p := repositories.ParamsForFetchUsers{
		CommonParamsForFetch: contexts.CommonParamsForFetch{
			Page:  uint64(req.Page),
			Limit: uint64(req.Limit),
		},
	}

	users, count, err = u.userRepo.Fetch(ctx, p)
	if err != nil {
		return
	}

	return
}

func (u *userUsecase) Update(ctx context.Context, id int64, req dto.UpdateUserReq) (err error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err = utils.CueValidateObject("UpdateUserRequest", cue.CueDefinitionForUser, req); err != nil {
		return err
	}

	copier.Copy(user, req)
	err = u.userRepo.Update(ctx, user, []string{"first_name", "last_name", "phone", "status"})
	return
}

func (u *userUsecase) Delete(ctx context.Context, id int64) (err error) {
	_, err = u.userRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	err = u.userRepo.DeleteById(ctx, id)
	return
}

func (u *userUsecase) GetByCode(ctx context.Context, code string) (user *domains.User, err error) {
	user, err = u.userRepo.GetByCode(ctx, code)
	return
}

func (u *userUsecase) GetByEmail(ctx context.Context, email string) (user *domains.User, err error) {
	user, err = u.userRepo.GetByEmail(ctx, email)
	return
}

func (u userUsecase) Register(ctx context.Context, req dto.CreateUserReq) (user *domains.User, err error) {
	user, err = u.Create(ctx, req)
	if err != nil {
		return
	}

	_, err = u.orgUsecase.Create(ctx, dto.CreateOrgReq{
		Name:   "My Organization",
		UserId: user.Id,
	})

	return
}
