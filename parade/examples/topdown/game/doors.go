package game

import "github.com/hajimehoshi/ebiten/v2"

const (
	DoorLength = 128
	DoorWidth  = 68
)

type Door struct {
	X, Y float64
}

func (d *Door) DrawDepth(dst *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Scale(CrateSize, CrateSize)
	geom.Translate(d.X, d.Y)
	dst.DrawImage(ImageWhite, &ebiten.DrawImageOptions{
		GeoM: geom,
	})
}

func (d *Door) DrawColor(dst *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Translate(d.X, d.Y)
	colorScale := ebiten.ColorScale{}
	colorScale.Scale(0.75, 0, 0.25, 1)
	colorScale.ScaleAlpha(0.25)
	dst.DrawImage(ImageCrateColor, &ebiten.DrawImageOptions{
		GeoM:       geom,
		ColorScale: colorScale,
	})
}
