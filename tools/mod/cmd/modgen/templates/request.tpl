package requests

// Target: delivery/requests/{{ .ModuleName | ToSnake }}.go

// Create{{ .ModuleName }}Req represent create {{ .ModuleName | ToLower }} request body
type Create{{ .ModuleName }}Req struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// Update{{ .ModuleName }}Req represent update {{ .ModuleName | ToLower }} request body
type Update{{ .ModuleName }}Req struct {
	ID    int64  `json:"-" param:"{{ .ModuleName | ToLowerCamel }}Id"`
	Name  string `json:"name"`
}

type Search{{ .ModuleName }}sReq struct {
	Limit int64 `json:"limit" query:"limit"`
    Page  int64 `json:"page" query:"page"`
}
