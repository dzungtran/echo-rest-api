package domains

import (
	"time"
)

type Project struct {
	Id              int64  `json:"id" db:"id"`
	Name            string `json:"name" db:"name"`
	OrgId           int64  `json:"org_id" db:"org_id"`
	Code            string `json:"code" db:"code"`
	Description     string `json:"description" db:"description"`
	Timezone        string `json:"timezone" db:"timezone"`
	Settings        string `json:"settings" db:"settings"`
	DefaultLanguage string `json:"default_language" db:"default_language"`

	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	ArchivedAt time.Time `json:"archived_at" db:"archived_at"`
}
