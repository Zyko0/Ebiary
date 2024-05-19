package internal

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var WhiteImage = ebiten.NewImage(3, 3)

func init() {
	WhiteImage.Fill(color.White)
}
