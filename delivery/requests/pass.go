package requests

// Target: delivery/requests/pass.go

// CreatePassReq represent create pass request body
type CreatePassReq struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// UpdatePassReq represent update pass request body
type UpdatePassReq struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type SearchPassesReq struct {
	Limit int64 `json:"limit" query:"limit"`
	Page  int64 `json:"page" query:"page"`
}
