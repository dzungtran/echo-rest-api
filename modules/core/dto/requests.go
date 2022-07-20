package dto

// CreateOrgReq represent create org request body
type CreateOrgReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Domain      string `json:"domain"`
	Logo        string `json:"logo"`
	UserId      int64  `json:"-"`
}

// UpdateOrgReq represent update org request body
type UpdateOrgReq struct {
	OrgId       int64  `json:"-" param:"orgId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Domain      string `json:"domain"`
	Logo        string `json:"logo"`
}

type SearchOrgsReq struct {
	Limit int64 `query:"limit"`
	Page  int64 `query:"page"`
	Ids   []int64
}

type InviteUsers struct {
	OrgId  int64    `json:"-" param:"orgId"`
	Emails []string `json:"emails"`
}

type CreateUserReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Code      string `json:"code"`
}

// UpdateUserReq represent update org request body
type UpdateUserReq struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
}

type SearchUsersReq struct {
	Limit int64 `json:"limit" query:"limit"`
	Page  int64 `json:"page" query:"page"`
}
