package domains

// Target: domains/{{ .ModuleName | ToSnake }}.go

import (
	"time"
)

type {{ .ModuleName }} struct {
	Id        int64      `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}
