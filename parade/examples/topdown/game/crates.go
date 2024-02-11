package game

import "github.com/hajimehoshi/ebiten/v2"

const (
	CrateSize = 64
)

type Crate struct {
	X, Y float64
}

func (c *Crate) DrawDepth(dst *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Scale(CrateSize, CrateSize)
	geom.Translate(c.X-CrateSize/2, c.Y-CrateSize/2)
	dst.DrawImage(ImageWhite, &ebiten.DrawImageOptions{
		GeoM: geom,
	})
}

func (c *Crate) DrawColor(dst *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Translate(c.X-CrateSize/2, c.Y-CrateSize/2)
	dst.DrawImage(ImageCrateColor, &ebiten.DrawImageOptions{
		GeoM: geom,
	})
}
