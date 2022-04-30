package definitions

import (
	"strings"
)

_StringRequired: string & !=""
_UserStatuses:   "active" | "deactivated" | "banned" | *"active"
_EmailRegex:     =~"(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$)"

// User info
#User: {
	id?: int & >=0
	// Generated Code, must be unique
	code:        _StringRequired & strings.MinRunes(10) & strings.MaxRunes(50)
	first_name?: string
	last_name?:  string
	email:       _StringRequired & strings.MinRunes(10) & _EmailRegex
	phone?:      string
	status:      _UserStatuses

	created_at?: string
	updated_at?: string
}

#CreateUserRequest: {
	first_name?: string
	last_name?:  string
	email:       _StringRequired & strings.MinRunes(10) & _EmailRegex
	phone?:      string
	code:        _StringRequired & strings.MinRunes(10)
}

#UpdateUserRequest: {
	first_name?: string
	last_name?:  string
	phone?:      string
	status:      _UserStatuses
}
