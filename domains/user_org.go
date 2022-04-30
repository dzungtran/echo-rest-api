package domains

import "time"

// Target: domains/user_org.go

type UserOrgRole string

const (
	UserRoleOwner   UserOrgRole = "owner"
	UserRoleManager UserOrgRole = "manager"
	UserRoleEditor  UserOrgRole = "editor"
	UserRoleViewer  UserOrgRole = "viewer"
	UserRoleGuest   UserOrgRole = "guest"
)

type UserOrg struct {
	Id        int64       `json:"id" db:"id"`
	UserId    int64       `json:"user_id" db:"user_id"`
	OrgId     int64       `json:"org_id" db:"org_id"`
	Role      UserOrgRole `json:"role" db:"role"`
	Status    UserStatus  `json:"status" db:"status"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt time.Time   `json:"updated_at" db:"updated_at"`
}

type UserWithRoles struct {
	User
	OrgRole map[int64]string `json:"org_role"`
}

func (u UserWithRoles) GetOrgIds() []int64 {
	ids := make([]int64, 0)
	if len(u.OrgRole) > 0 {
		for id := range u.OrgRole {
			ids = append(ids, id)
		}
	}
	return ids
}
