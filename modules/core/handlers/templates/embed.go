package templates

import "embed"

//go:embed *.tpl
var CoreTemplates embed.FS
