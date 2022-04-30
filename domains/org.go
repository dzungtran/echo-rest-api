package domains

// Target: domains/org.go

import (
	"time"

	"github.com/dzungtran/echo-rest-api/pkg/cue"
	"github.com/dzungtran/echo-rest-api/pkg/utils"
)

type Org struct {
	Id          int64  `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Code        string `json:"code" db:"code"`
	Description string `json:"description" db:"description"`
	Domain      string `json:"domain" db:"domain"`
	Logo        string `json:"logo" db:"logo"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func (u Org) Verify() error {
	return utils.CueValidateObject("Org", cue.CueDefinitionForUser, u)
}
