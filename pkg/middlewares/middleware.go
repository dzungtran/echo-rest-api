package middlewares

import (
	"context"
	"net/http"
	"reflect"
	"strings"

	"github.com/dzungtran/echo-rest-api/config"
	"github.com/dzungtran/echo-rest-api/modules/core/domains"
	coreRepo "github.com/dzungtran/echo-rest-api/modules/core/repositories"
	projectRepo "github.com/dzungtran/echo-rest-api/modules/projects/repositories"
	"github.com/dzungtran/echo-rest-api/pkg/authz"
	"github.com/dzungtran/echo-rest-api/pkg/constants"
	"github.com/dzungtran/echo-rest-api/pkg/contexts"
	"github.com/dzungtran/echo-rest-api/pkg/kratos"
	"github.com/dzungtran/echo-rest-api/pkg/logger"
	"github.com/dzungtran/echo-rest-api/pkg/utils"
	"github.com/labstack/echo/v4"
	ory "github.com/ory/kratos-client-go"
)

// MiddlewareManager ...
// This file contains common functions for auth
type MiddlewareManager struct {
	appConf      *config.AppConfig
	userRepo     coreRepo.UserRepository
	userOrgRepo  coreRepo.UserOrgRepository
	projectRepo  projectRepo.ProjectRepository
	orgRepo      coreRepo.OrgRepository
	kratosClient *ory.APIClient
}

// NewMiddlewareManager will create new an MiddlewareManager object
func NewMiddlewareManager(
	appConf *config.AppConfig,
	userRepo coreRepo.UserRepository,
	userOrgRepo coreRepo.UserOrgRepository,
	projectRepo projectRepo.ProjectRepository,
	orgRepo coreRepo.OrgRepository,
) *MiddlewareManager {
	return &MiddlewareManager{
		appConf:      appConf,
		userRepo:     userRepo,
		userOrgRepo:  userOrgRepo,
		projectRepo:  projectRepo,
		orgRepo:      orgRepo,
		kratosClient: kratos.NewKratosSelfHostedClient(appConf.KratosApiEndpoint, appConf.Environment == "development"),
	}
}

func (m MiddlewareManager) Auth() echo.MiddlewareFunc {
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
				logger.Log().Errorw("error while fetch user for auth", "email", strings.ReplaceAll(email, "\n", ""), "error", err)
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": "cannot fetch user",
				})
			}

			if u.Status != domains.UserStatusActive {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{
					"error": "user is not active",
				})
			}

			c.Set(constants.ContextKeyUser, u)
			return next(c)
		}
	}
}

func (m MiddlewareManager) CheckPolicies() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			denyMsg, err := authz.CheckPoliciesContext(c)
			if err != nil {
				msg := ""
				if len(denyMsg) > 0 {
					msg = denyMsg[0]
				}
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"error": msg,
				})
			}
			return next(c)
		}
	}
}

func (m MiddlewareManager) CheckPoliciesWithOrg() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			orgId := utils.GetResourceIdFromParam(c, "orgId")
			org, err := m.orgRepo.GetByID(c.Request().Context(), orgId)
			if err != nil {
				if err == constants.ErrNotFound {
					return c.JSON(http.StatusNotFound, map[string]interface{}{
						"error": "Not found",
					})
				}
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error": err.Error(),
				})
			}

			denyMsg, err := authz.CheckPoliciesContext(c, authz.WithInputOrg(org))
			if err != nil {
				msg := ""
				if len(denyMsg) > 0 {
					msg = denyMsg[0]
				}
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"error": msg,
				})
			}

			c.Set(constants.ContextKeyOrg, org)
			return next(c)
		}
	}
}

func (m MiddlewareManager) CheckPoliciesWithProject() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			projectId := utils.GetResourceIdFromParam(c, "projectId")
			project, err := m.projectRepo.GetByID(c.Request().Context(), projectId)
			if err != nil {
				if err == constants.ErrNotFound {
					return c.JSON(http.StatusNotFound, map[string]interface{}{
						"error": "Not found",
					})
				}
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error": err.Error(),
				})
			}

			denyMsg, err := authz.CheckPoliciesContext(c, authz.WithInputExtraData("project", project))
			if err != nil {
				msg := ""
				if len(denyMsg) > 0 {
					msg = denyMsg[0]
				}
				return c.JSON(http.StatusForbidden, map[string]interface{}{
					"error": msg,
				})
			}

			c.Set(constants.ContextKeyProject, project)
			return next(c)
		}
	}
}

func (m MiddlewareManager) CheckPoliciesWithRequestPayload(payloadInst interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			if reflect.ValueOf(payloadInst).Kind() != reflect.Ptr {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{
					"error": "payload should be a pointer",
				})
			}

			if err := c.Bind(payloadInst); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"error": err.Error(),
				})
			}

			c.Set(constants.ContextKeyPayload, payloadInst)
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

	userOrgs, _, err := m.userOrgRepo.Fetch(ctx, coreRepo.ParamsForFetchUserOrgs{
		CommonParamsForFetch: contexts.CommonParamsForFetch{
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
