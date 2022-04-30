package deny

import data.utils

deny_user_endpoint := deny_get_user 
| deny_update_user 
| deny_delete_user 
| deny_get_list_user 
| deny_create_user

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

# START Get user info
deny_get_user[msg] {
	# invalid user info
	is_get_user_info
	not input.user_info
	msg := "user id is invalid"
}

deny_get_user[msg] {
	# invalid user info
	is_get_user_info
	not utils.is_super_admin
	input.user_info.id <= 0
	msg := "user id is invalid"
}

deny_get_user[msg] {
	# invalid user info
	is_get_user_info
	not input.user_info.id
	msg := "user id is invalid"
}

deny_get_user[msg] {
	is_get_user_info
	input.user.id != input.user_info.id
	msg := "user id is invalid"
}

# END Get user info

# START Update user info
deny_update_user[msg] {
	is_update_user_info
	not input.user_info
	msg := "user id is invalid"
}

deny_update_user[msg] {
	is_update_user_info
	input.user_info.id <= 0
	msg := "user id is invalid"
}

deny_update_user[msg] {
	# invalid user info
	is_update_user_info
	not input.user_info.id
	msg := "user id is invalid"
}

deny_update_user[msg] {
	is_update_user_info
	input.user.id != input.user_info.id
	msg := "user id is invalid"
}

# END Update user info

deny_delete_user[msg] {
	input.method == "DELETE"
	input.endpoint == "/admin/users/:userId"
	not utils.is_super_admin
	msg := "you don't have the permission"
}

deny_get_list_user[msg] {
	input.method == "GET"
	input.endpoint == "/admin/users"
	not utils.is_super_admin
	msg := "you don't have the permission"
}

deny_create_user[msg] {
	input.method == "POST"
	input.endpoint == "/admin/users"
	not utils.is_super_admin
	msg := "you don't have the permission"
}