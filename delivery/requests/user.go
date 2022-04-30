package requests

// CreateUserReq represent create org request body
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
