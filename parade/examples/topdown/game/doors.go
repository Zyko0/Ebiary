package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	DoorLength          = 128
	DoorWidth           = 10
	DoorTriggerDistance = 256

	doorDepthSpeed = 0.025
)

type Door struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
	Depth  float64

	depthDirection float64
}

var (
	Doors = []*Door{
		// West
		{
			depthDirection: 1,
			X:              68 - DoorWidth/2 - 1,
			Y:              1024 / 2,
			Width:          DoorWidth,
			Height:         DoorLength,
			Depth:          1,
		},
		// East
		{
			depthDirection: 1,
			X:              1024 - 68 + DoorWidth/2 + 1,
			Y:              1024 / 2,
			Width:          DoorWidth,
			Height:         DoorLength,
			Depth:          1,
		},
		// North
		{
			depthDirection: 1,
			X:              1024 / 2,
			Y:              68 - DoorWidth/2 - 1,
			Width:          DoorLength,
			Height:         DoorWidth,
			Depth:          1,
		},
		// South
		{
			depthDirection: 1,
			X:              1024 / 2,
			Y:              1024 - 68 + DoorWidth/2 + 1,
			Width:          DoorLength,
			Height:         DoorWidth,
			Depth:          1,
		},
	}
)

func (d *Door) DistanceTo(p *Player) float64 {
	dx := math.Abs(p.X-d.X) - d.Width
	dy := math.Abs(p.Y-d.Y) - d.Height
	adx := max(dx, 0)
	ady := max(dy, 0)

	return math.Sqrt(adx*adx+ady*ady) + min(max(dx, dy), 0)
}

func (d *Door) TriggerActivation(opening bool) {
	if opening {
		d.depthDirection = -1
	} else {
		d.depthDirection = 1
	}
}

func (d *Door) Update() {
	d.Depth += d.depthDirection * doorDepthSpeed
	d.Depth = max(min(d.Depth, 1), 0)
}

func (d *Door) DrawDepth(dst *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Scale(d.Width, d.Height)
	geom.Translate(d.X-d.Width/2, d.Y-d.Height/2)
	cs := ebiten.ColorScale{}
	h := float32(d.Depth)
	cs.Scale(h, h, h, 1)
	dst.DrawImage(ImageWhite, &ebiten.DrawImageOptions{
		GeoM:       geom,
		ColorScale: cs,
	})
}

func (d *Door) DrawColor(dst *ebiten.Image) {
	geom := ebiten.GeoM{}
	// Little hack, drawing more color on the texture for boxmapping
	size := max(d.Width, d.Height)
	geom.Scale(size, size)
	geom.Translate(d.X-size/2, d.Y-size/2)
	cs := ebiten.ColorScale{}
	cs.Scale(1, 0.2, 0.4, 1)
	cs.ScaleAlpha(0.75)
	dst.DrawImage(ImageWhite, &ebiten.DrawImageOptions{
		GeoM:       geom,
		ColorScale: cs,
	})
}
