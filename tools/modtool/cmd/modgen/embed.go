package modgen

import "embed"

//go:embed templates/*.tpl
var fsTemplates embed.FS

func GetModGenTemplates() embed.FS {
	return fsTemplates
}
