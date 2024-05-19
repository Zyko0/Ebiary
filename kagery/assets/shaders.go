package assets

import (
	_ "embed"
)

var (
	//go:embed shaders/example.kage
	ShaderSrc string
)