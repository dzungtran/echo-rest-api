package dto

// Create{{ .SingularName }}Req represent create {{ .SingularName }} request body
type Create{{ .SingularName }}Req struct {
	Name string `json:"name"`
}

// Update{{ .SingularName }}Req represent update {{ .SingularName }} request body
type Update{{ .SingularName }}Req struct {
	ID   int64  `json:"-" param:"{{ .SingularName | ToLowerCamel }}Id"`
	Name string `json:"name"`
}

type Search{{ .PluralName }}Req struct {
	Limit int64 `json:"limit" query:"limit"`
	Page  int64 `json:"page" query:"page"`
}
