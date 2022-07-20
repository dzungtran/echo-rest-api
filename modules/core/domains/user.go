package domains

import (
	"time"

	"github.com/dzungtran/echo-rest-api/pkg/cue"
	"github.com/dzungtran/echo-rest-api/pkg/utils"
)

type UserStatus string

const (
	UserStatusActive      UserStatus = "active"
	UserStatusDeactivated UserStatus = "deactivated"
	UserStatusBanned      UserStatus = "banned"

	// user org
	UserStatusInvited UserStatus = "invited"
)

type User struct {
	Id        int64      `json:"id" db:"id"`
	FirstName string     `json:"first_name" db:"first_name"`
	LastName  string     `json:"last_name" db:"last_name"`
	Code      string     `json:"code" db:"code"`
	Email     string     `json:"email" db:"email"`
	Phone     string     `json:"phone" db:"phone"`
	Status    UserStatus `json:"status" db:"status"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

func (u User) Verify() error {
	return utils.CueValidateObject("User", cue.CueDefinitionForUser, u)
}
