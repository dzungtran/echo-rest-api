package authz

import (
	"testing"

	"github.com/dzungtran/echo-rest-api/modules/core/domains"
	"github.com/stretchr/testify/assert"
)

var (
	getOrgEndpoint     = TestEndpoint{"GET", "/admin/orgs/:orgId"}
	updateOrgEndpoint  = TestEndpoint{"PUT", "/admin/orgs/:orgId"}
	deleteOrgEndpoint  = TestEndpoint{"DELETE", "/admin/orgs/:orgId"}
	getListOrgEndpoint = TestEndpoint{"GET", "/admin/orgs"}
	createOrgEndpoint  = TestEndpoint{"POST", "/admin/orgs"}
)

func TestPoliciesForOrgEndpoint(t *testing.T) {
	loggedInUser := &domains.UserWithRoles{
		User: domains.User{Id: 8},
		OrgRole: map[int64]string{
			7:  "guest",
			8:  "viewer",
			9:  "manager",
			10: "owner",
		},
	}

	superAdmin := &domains.UserWithRoles{
		User:    domains.User{Id: 99, Email: "hello@iamdzung.com"},
		OrgRole: nil,
	}

	tcs := []struct {
		name         string
		loggedInUser *domains.UserWithRoles
		requestedOrg *domains.Org
		hasError     bool
		denyMsg      []string
		endpoint     TestEndpoint
	}{
		{
			"should allow super admin can fetch org info of others",
			superAdmin,
			&domains.Org{Id: 8},
			false,
			[]string{},
			getOrgEndpoint,
		},
		{
			"should allow viewer can get org info",
			loggedInUser,
			&domains.Org{Id: 8},
			false,
			[]string{},
			getOrgEndpoint,
		},
		{
			"should deny get org info with user not in org",
			loggedInUser,
			&domains.Org{Id: 99},
			true,
			[]string{},
			getOrgEndpoint,
		},
		{
			"should allow any logged in user can create an org",
			loggedInUser,
			&domains.Org{},
			false,
			[]string{},
			createOrgEndpoint,
		},
		{
			"should allow owners can update their orgs",
			loggedInUser,
			&domains.Org{Id: 10},
			false,
			[]string{},
			updateOrgEndpoint,
		},
		{
			"should deny other roles try update org info",
			loggedInUser,
			&domains.Org{Id: 9},
			true,
			[]string{},
			updateOrgEndpoint,
		},
		{
			"should allow owners can delete their orgs",
			loggedInUser,
			&domains.Org{Id: 10},
			false,
			[]string{},
			deleteOrgEndpoint,
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			msg, err := CheckPolicies(tc.loggedInUser,
				WithInputRequestMethod(tc.endpoint.Method),
				WithInputRequestEndpoint(tc.endpoint.Endpoint),
				WithInputOrg(tc.requestedOrg),
			)
			if tc.hasError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}

			if len(tc.denyMsg) > 0 {
				assert.Equal(t, tc.denyMsg, msg)
			}

			assert.Equal(t, len(tc.denyMsg), len(msg))
		})
	}
}
