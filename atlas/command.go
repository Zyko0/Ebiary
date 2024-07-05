package atlas

import "github.com/hajimehoshi/ebiten/v2"

type DrawCommand struct {
	Image *Image

	ColorScale ebiten.ColorScale
	GeoM       ebiten.GeoM
}

type DrawOptions struct {
	Blend  ebiten.Blend
	Filter ebiten.Filter
}

func Draw(dst *ebiten.Image, opts *DrawOptions, commands ...*DrawCommand) {
	
}
