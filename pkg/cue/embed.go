package cue

import (
	_ "embed"
)

var (
	//go:embed definitions/org.cue
	CueDefinitionForOrg string

	//go:embed definitions/user.cue
	CueDefinitionForUser string

	//go:embed definitions/project.cue
	CueDefinitionForProject string
)
