package cue

import (
	_ "embed"
)

var (
	//go:embed defs/project.cue
	CueDefinitionForProject string
)
