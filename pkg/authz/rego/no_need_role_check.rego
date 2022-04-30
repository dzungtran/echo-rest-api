package authz

import data.utils

default no_need_role_check = []

no_need_role_check = no_need_role_check_org_endpoint | no_need_role_check_user_endpoint

is_get_user_info {
	not utils.is_super_admin
	input.method == "GET"
	input.endpoint == "/admin/users/:userId"
}

is_update_user_info {
	not utils.is_super_admin
	input.method == "PUT"
	input.endpoint == "/admin/users/:userId"
}

no_need_role_check_user_endpoint[act] {
	is_update_user_info
	act := {
		"endpoint": input.endpoint,
		"method": input.method,
	}
}

no_need_role_check_user_endpoint[act] {
	is_get_user_info
	act := {
		"endpoint": input.endpoint,
		"method": input.method,
	}
}

no_need_role_check_org_endpoint[act] {
	# Get list org
	input.endpoint = "/admin/orgs"
	input.method = "GET"
	act = {
		"endpoint": input.endpoint,
		"method": input.method,
	}
}

no_need_role_check_org_endpoint[act] {
	# Get create an org
	input.endpoint = "/admin/orgs"
	input.method = "POST"
	act = {
		"endpoint": input.endpoint,
		"method": input.method,
	}
}
