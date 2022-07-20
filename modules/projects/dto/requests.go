package dto

// CreateProjectReq represent create project request body
type CreateProjectReq struct {
	OrgId           int64  `json:"org_id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Timezone        string `json:"timezone"`
	DefaultLanguage string `json:"default_language"`
}

// UpdateProjectReq represent update project request body
type UpdateProjectReq struct {
	ID              int64  `json:"-" param:"projectId"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	Timezone        string `json:"timezone"`
	DefaultLanguage string `json:"default_language"`
}

type SearchProjectsReq struct {
	Limit int64 `json:"limit" query:"limit"`
	Page  int64 `json:"page" query:"page"`
	OrgId int64 `json:"org_id" query:"org_id"`
}
