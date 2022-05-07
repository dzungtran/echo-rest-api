package utils

import future.keywords.in

default is_super_admin = false
is_super_admin {
	input.user.email in [
		"hello@iamdzung.com",
	]
}
