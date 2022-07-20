package authz

import (
	"testing"

	"github.com/dzungtran/echo-rest-api/modules/core/domains"
	"github.com/stretchr/testify/assert"
)

type TestEndpoint struct {
	Method   string
	Endpoint string
}

var (
	getUserEndpoint     = TestEndpoint{"GET", "/admin/users/:userId"}
	updateUserEndpoint  = TestEndpoint{"PUT", "/admin/users/:userId"}
	deleteUserEndpoint  = TestEndpoint{"DELETE", "/admin/users/:userId"}
	getListUserEndpoint = TestEndpoint{"GET", "/admin/users"}
	createUserEndpoint  = TestEndpoint{"POST", "/admin/users"}
)

func TestPoliciesForUserEndpoint(t *testing.T) {
	loggedInUser := &domains.UserWithRoles{
		User: domains.User{Id: 8},
		OrgRole: map[int64]string{
			8: "viewer",
		},
	}

	superAdmin := &domains.UserWithRoles{
		User:    domains.User{Id: 99, Email: "hello@iamdzung.com"},
		OrgRole: map[int64]string{},
	}

	tcs := []struct {
		name          string
		loggedInUser  *domains.UserWithRoles
		requestedUser *domains.User
		hasError      bool
		denyMsg       []string
		endpoint      TestEndpoint
	}{
		{
			"super admin fetch user info then return no error",
			superAdmin,
			&domains.User{Id: 8},
			false,
			[]string{},
			getUserEndpoint,
		},
		{
			"logged in users fetch info by them self then return no error",
			loggedInUser,
			&domains.User{Id: 8},
			false,
			[]string{},
			getUserEndpoint,
		},
		{
			"logged in users fetch info of other then return forbidden error",
			loggedInUser,
			&domains.User{Id: 9},
			true,
			[]string{"user id is invalid"},
			getUserEndpoint,
		},
		{
			"logged in users fetch null info then return forbidden error",
			loggedInUser,
			nil,
			true,
			[]string{"user id is invalid"},
			getUserEndpoint,
		},

		{
			"super admin update user info then return no error",
			superAdmin,
			&domains.User{Id: 8},
			false,
			[]string{},
			updateUserEndpoint,
		},
		{
			"logged in users update info by them self then return no error",
			loggedInUser,
			&domains.User{Id: 8},
			false,
			[]string{},
			updateUserEndpoint,
		},
		{
			"logged in users update info of other then return forbidden error",
			loggedInUser,
			&domains.User{Id: 9},
			true,
			[]string{"user id is invalid"},
			updateUserEndpoint,
		},
		{
			"logged in users update null info then return forbidden error",
			loggedInUser,
			nil,
			true,
			[]string{"user id is invalid"},
			updateUserEndpoint,
		},

		{
			"super admin delete an user then return no error",
			superAdmin,
			&domains.User{Id: 99},
			false,
			[]string{},
			deleteUserEndpoint,
		},
		{
			"logged in users try to delete an user then return forbidden error",
			loggedInUser,
			&domains.User{Id: 99},
			true,
			[]string{"you don't have the permission"},
			deleteUserEndpoint,
		},

		{
			"super admin create an user then return no error",
			superAdmin,
			&domains.User{Id: 99},
			false,
			[]string{},
			createUserEndpoint,
		},
		{
			"logged in users try to create an user then return forbidden error",
			loggedInUser,
			&domains.User{Id: 99},
			true,
			[]string{"you don't have the permission"},
			createUserEndpoint,
		},

		{
			"super admin get list user then return no error",
			superAdmin,
			&domains.User{Id: 99},
			false,
			[]string{},
			getListUserEndpoint,
		},
		{
			"logged in users try to get list user then return forbidden error",
			loggedInUser,
			&domains.User{Id: 99},
			true,
			[]string{"you don't have the permission"},
			getListUserEndpoint,
		},
	}

	for _, tc := range tcs {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			msg, err := CheckPolicies(tc.loggedInUser,
				WithInputRequestMethod(tc.endpoint.Method),
				WithInputRequestEndpoint(tc.endpoint.Endpoint),
				WithInputExtraData("user_info", tc.requestedUser),
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
