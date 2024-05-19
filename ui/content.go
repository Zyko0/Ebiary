package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Content interface {
	Draw(dst *ebiten.Image)
}

