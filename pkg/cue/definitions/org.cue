package definitions

import (
	"strings"
)

_ASCIIChars:  string & =~"^[\\x00-\\x7F]+$"
_OrgStatuses: "active" | "inactive"
_EmailRegex:  =~"(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$)"

// Org info
#Org: {
	id?: int & >=0
	// Org name
	name: _ASCIIChars & !="" & strings.MinRunes(10) & strings.MaxRunes(100)
	// Generated Code
	code:         string & !="" & strings.MinRunes(10) & strings.MaxRunes(50)
	description?: string
	domain?:      string
	logo?:        string
	status:       _OrgStatuses

	created_at?: string
	updated_at?: string
}

#CreateOrgRequest: {
	// Org name
	name:         _ASCIIChars & !="" & strings.MinRunes(10) & strings.MaxRunes(100)
	description?: string
	domain?:      string
	logo?:        string
}

#UpdateOrgRequest: {
	// Org name
	name:         _ASCIIChars & !="" & strings.MinRunes(10) & strings.MaxRunes(100)
	description?: string
	domain?:      string
	logo?:        string
	status:       _OrgStatuses
}

#InviteOrgRequest: {
	emails: [..._EmailRegex]
}
