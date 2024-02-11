package game

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	PlayerSize          = 128
	PlayerHeight        = 64
	PlayerMovementSpeed = 5
)

type Player struct {
	X, Y        float64
	Orientation float64
}

func NewPlayer() *Player {
	return &Player{
		X: 128,
		Y: 128, // TODO:
	}
}

func (p *Player) Update() {
	var dx, dy float64

	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		dx = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		dx = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		dy = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		dy = 1
	}
	// Update orientation
	if dx != 0 || dy != 0 {
		p.Orientation = math.Atan2(dx, -dy)
	}
	// Nerf diagonal movement speed
	if dx != 0 && dy != 0 {
		const diagonalFactor = math.Sqrt2 / 2
		dx *= diagonalFactor
		dy *= diagonalFactor
	}
	// Update player coordinates
	p.X += dx * PlayerMovementSpeed
	p.Y += dy * PlayerMovementSpeed
}

func (p *Player) DrawDepth(dst *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Translate(-PlayerSize/2, -PlayerSize/2)
	geom.Rotate(p.Orientation)
	geom.Translate(p.X, p.Y)
	dst.DrawImage(ImagePlayerDepth, &ebiten.DrawImageOptions{
		GeoM: geom,
	})
}

func (p *Player) DrawColor(dst *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Translate(-PlayerSize/2, -PlayerSize/2)
	geom.Rotate(p.Orientation)
	geom.Translate(p.X, p.Y)
	dst.DrawImage(ImagePlayerColor, &ebiten.DrawImageOptions{
		GeoM: geom,
	})
}
