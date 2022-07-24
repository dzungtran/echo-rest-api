package usecases

// Target: usecases/org.go

import (
	"context"
	"errors"

	"github.com/dzungtran/echo-rest-api/modules/core/domains"
	"github.com/dzungtran/echo-rest-api/modules/core/dto"
	"github.com/dzungtran/echo-rest-api/modules/core/repositories"
	"github.com/dzungtran/echo-rest-api/pkg/contexts"
	"github.com/dzungtran/echo-rest-api/pkg/cue"
	sqlTools "github.com/dzungtran/echo-rest-api/pkg/sql-tools"
	"github.com/dzungtran/echo-rest-api/pkg/utils"
	"github.com/jinzhu/copier"
)

// OrgUsecase represent the org's usecase contract
type OrgUsecase interface {
	Create(ctx context.Context, request dto.CreateOrgReq) (*domains.Org, error)
	GetByID(ctx context.Context, id int64) (*domains.Org, error)
	Fetch(ctx context.Context, req dto.SearchOrgsReq) ([]*domains.Org, int64, error)
	Update(ctx context.Context, id int64, request dto.UpdateOrgReq) error
	Delete(ctx context.Context, id int64) error
	Invite(ctx context.Context, id int64, request dto.InviteUsers) error
}

type orgUsecase struct {
	orgRepo     repositories.OrgRepository
	userOrgRepo repositories.UserOrgRepository
	sqlxTrans   *sqlTools.SqlxTransaction
}

// NewOrgUsecase will create new an orgUsecase object representation of OrgUsecase interface
func NewOrgUsecase(
	orgRepo repositories.OrgRepository,
	userOrgRepo repositories.UserOrgRepository,
	sqlxTrans *sqlTools.SqlxTransaction,
) OrgUsecase {
	return &orgUsecase{
		orgRepo:     orgRepo,
		userOrgRepo: userOrgRepo,
		sqlxTrans:   sqlxTrans,
	}
}

func (u *orgUsecase) Create(ctx context.Context, req dto.CreateOrgReq) (org *domains.Org, err error) {
	org = &domains.Org{}
	var orgId int64

	if req.UserId < 0 {
		err = errors.New("missing user")
		return
	}

	if err = utils.CueValidateObject("CreateOrgRequest", cue.CueDefinitionForOrg, req); err != nil {
		return nil, err
	}

	copier.Copy(org, req)
	org.Code = utils.GenerateLongUUID()

	// Start transaction
	tx, err := u.sqlxTrans.Init()
	if err != nil {
		return
	}

	needRollback := false
	defer func() {
		if needRollback {
			tx.Rollback()
		}
	}()

	orgId, err = u.orgRepo.CreateWithTx(ctx, tx, org)
	if err != nil {
		needRollback = true
		return
	}

	_, err = u.userOrgRepo.CreateWithTx(ctx, tx, &domains.UserOrg{
		UserId: req.UserId,
		OrgId:  orgId,
		Role:   domains.UserRoleOwner,
		Status: domains.UserStatusActive,
	})
	if err != nil {
		needRollback = true
		return
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return
	}
	// END transction

	org, err = u.orgRepo.GetByID(ctx, orgId)
	return
}

func (u *orgUsecase) GetByID(ctx context.Context, id int64) (org *domains.Org, err error) {
	org, err = u.orgRepo.GetByID(ctx, id)
	return
}

func (u *orgUsecase) Fetch(ctx context.Context, req dto.SearchOrgsReq) (orgs []*domains.Org, count int64, err error) {
	p := repositories.ParamsForFetchOrgs{
		Ids: req.Ids,
		CommonParamsForFetch: contexts.CommonParamsForFetch{
			Page:  uint64(req.Page),
			Limit: uint64(req.Limit),
		},
	}

	orgs, count, err = u.orgRepo.Fetch(ctx, p)
	if err != nil {
		return
	}

	return
}

func (u *orgUsecase) Update(ctx context.Context, id int64, req dto.UpdateOrgReq) (err error) {
	org, err := u.orgRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if err = utils.CueValidateObject("UpdateOrgRequest", cue.CueDefinitionForOrg, req); err != nil {
		return err
	}

	copier.Copy(org, req)
	err = u.orgRepo.Update(ctx, org, []string{"name", "description", "logo", "domain"})
	return
}

func (u *orgUsecase) Delete(ctx context.Context, id int64) (err error) {
	_, err = u.orgRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	err = u.orgRepo.DeleteById(ctx, id)
	return
}

func (u *orgUsecase) Invite(ctx context.Context, orgId int64, req dto.InviteUsers) (err error) {
	_, err = u.orgRepo.GetByID(ctx, orgId)
	if err != nil {
		return err
	}

	if len(req.Emails) == 0 || len(req.Emails) > 20 {
		return errors.New("can invites 1-20 emails at once")
	}

	if err = utils.CueValidateObject("InviteOrgRequest", cue.CueDefinitionForUser, req); err != nil {
		return err
	}

	// TODO: Create none existed users and user org here

	// TODO: Send email here

	return
}
