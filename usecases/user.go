package usecases

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/dzungtran/echo-rest-api/delivery/defines"
	"github.com/dzungtran/echo-rest-api/delivery/requests"
	"github.com/dzungtran/echo-rest-api/domains"
	"github.com/dzungtran/echo-rest-api/pkg/cue"
	"github.com/dzungtran/echo-rest-api/pkg/utils"
	"github.com/dzungtran/echo-rest-api/repositories/postgres"
)

// UserUsecase represent the user's usecase contract
type UserUsecase interface {
	Create(ctx context.Context, request requests.CreateUserReq) (*domains.User, error)
	GetByID(ctx context.Context, id int64) (*domains.User, error)
	Fetch(ctx context.Context, req requests.SearchUsersReq) ([]*domains.User, int64, error)
	Update(ctx context.Context, id int64, request requests.UpdateUserReq) error
	Delete(ctx context.Context, id int64) error
	GetByCode(ctx context.Context, code string) (*domains.User, error)
	GetByEmail(ctx context.Context, email string) (*domains.User, error)
}

type userUsecase struct {
	userRepo postgres.UserRepository
}

// NewUserUsecase will create new an userUsecase object representation of UserUsecase interface
func NewUserUsecase(userRepo postgres.UserRepository) UserUsecase {
	return &userUsecase{
		userRepo: userRepo,
	}
}

func (u *userUsecase) Create(ctx context.Context, req requests.CreateUserReq) (user *domains.User, err error) {
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

func (u *userUsecase) Fetch(ctx context.Context, req requests.SearchUsersReq) (users []*domains.User, count int64, err error) {
	p := postgres.ParamsForFetchUsers{
		CommonParamsForFetch: defines.CommonParamsForFetch{
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

func (u *userUsecase) Update(ctx context.Context, id int64, req requests.UpdateUserReq) (err error) {
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
