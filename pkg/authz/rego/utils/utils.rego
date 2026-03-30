package utils

import future.keywords.if
import future.keywords.in

default is_super_admin := false

is_super_admin if {
	input.user.email in [
		"your_admin@email.com",
	]
}
