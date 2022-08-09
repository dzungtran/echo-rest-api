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

// User domain info
// @Description User account information
type User struct {
	// User indetifier number
	Id int64 `json:"id" db:"id" example:"1"`
	// User first name
	FirstName string `json:"first_name" db:"first_name" example:"Dzung"`
	// User last name
	LastName string `json:"last_name" db:"last_name" example:"Tran"`
	// User last name
	Code      string     `json:"code" db:"code" example:"95a8d1aa-xxx-xxx-0c15d41"`
	Email     string     `json:"email" db:"email" example:"email@api.com"`
	Phone     string     `json:"phone" db:"phone" example:"+84 0986415xxxx"`
	Status    UserStatus `json:"status" db:"status" example:"active" enums:"active,deactivated,banned"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

func (u User) Verify() error {
	return utils.CueValidateObject("User", cue.CueDefinitionForUser, u)
}
