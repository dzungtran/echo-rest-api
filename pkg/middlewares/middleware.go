package middlewares

import (
	"context"
	"net/http"

	"github.com/dzungtran/echo-rest-api/config"
	"github.com/dzungtran/echo-rest-api/delivery/defines"
	"github.com/dzungtran/echo-rest-api/domains"
	"github.com/dzungtran/echo-rest-api/pkg/constants"
	"github.com/dzungtran/echo-rest-api/pkg/kratos"
	"github.com/dzungtran/echo-rest-api/pkg/logger"
	"github.com/dzungtran/echo-rest-api/pkg/utils"
	"github.com/dzungtran/echo-rest-api/repositories/postgres"
	"github.com/labstack/echo/v4"
	ory "github.com/ory/kratos-client-go"
)

// MiddlewareManager ...
// This file contains common functions for auth
type MiddlewareManager struct {
	appConf      *config.AppConfig
	userRepo     postgres.UserRepository
	userOrgRepo  postgres.UserOrgRepository
	kratosClient *ory.APIClient
}

// NewMiddlewareManager will create new an MiddlewareManager object
func NewMiddlewareManager(
	appConf *config.AppConfig,
	userRepo postgres.UserRepository,
	userOrgRepo postgres.UserOrgRepository,
) *MiddlewareManager {
	return &MiddlewareManager{
		appConf:      appConf,
		userRepo:     userRepo,
		userOrgRepo:  userOrgRepo,
		kratosClient: kratos.NewKratosSelfHostedClient(appConf.KratosApiEndpoint, appConf.Environment == "development"),
	}
}

func (m *MiddlewareManager) GenerateRequestID() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			reqId := c.Request().Header.Get(constants.HeaderXRequestID)
			if reqId == "" {
				reqId = utils.GenerateUUID()
			}

			c.Request().Header.Set(constants.HeaderXRequestID, reqId)
			c.Set(constants.RequestIDContextKey, reqId)
			return next(c)
		}
	}
}

func (m *MiddlewareManager) Auth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		if m.appConf.Environment != "local" {
			// current just support Kratos Authn
			return m.KratosAuth(next)
		}

		// Default auth here for local debug and development
		// Bypass auth
		return func(c echo.Context) error {
			ctx := c.Request().Context()
			email := c.Request().Header.Get(constants.HeaderXUserEmail)

			u, err := m.fetchUserFromAuth(ctx, "", email)
			if err != nil {
				logger.Log().Errorw("error while fetch user for auth", "email", email, "error", err)
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"message": "cannot fetch user",
				})
			}

			if u.Status != domains.UserStatusActive {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"message": "user is not active",
				})
			}

			c.Set(constants.UserContextKey, u)
			return next(c)
		}
	}
}

func (m *MiddlewareManager) fetchUserFromAuth(ctx context.Context, code, email string) (u *domains.UserWithRoles, err error) {
	var user *domains.User
	if code != "" {
		user, err = m.userRepo.GetByCode(ctx, code)
		if err != nil {
			return
		}
	} else if email != "" {
		user, err = m.userRepo.GetByEmail(ctx, email)
		if err != nil {
			return
		}
	} else {
		return nil, constants.ErrUnauthorized
	}

	if user == nil {
		return nil, constants.ErrUnauthorized
	}

	u = &domains.UserWithRoles{
		User:    *user,
		OrgRole: map[int64]string{},
	}

	userOrgs, _, err := m.userOrgRepo.Fetch(ctx, postgres.ParamsForFetchUserOrgs{
		CommonParamsForFetch: defines.CommonParamsForFetch{
			NoLimit: true,
		},
		UserIds: []int64{user.Id},
	})
	if err != nil {
		return
	}

	for _, uo := range userOrgs {
		u.OrgRole[uo.OrgId] = string(uo.Role)
	}
	return
}
