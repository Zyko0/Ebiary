package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	PlayerSize    = 256
	FrameDuration = 20

	MovementSpeed   = 5
	MaxJumpVelocity = 25
	MaxFallVelocity = -20
)

const (
	StateIdle byte = iota
	StateRunning
	StateJumping
	StateFalling
)

type Player struct {
	State     byte
	StateTick int

	Grounded  bool
	DirX      int
	VelocityY float64
	X, Y      float64
}

func NewPlayer() *Player {
	return &Player{
		State: StateIdle,

		Grounded: true,
		DirX:     1,
		X:        0,
		Y:        0,
	}
}

func (p *Player) Update() {
	// Direction
	moveIntent := true
	switch {
	case ebiten.IsKeyPressed(ebiten.KeyLeft), ebiten.IsKeyPressed(ebiten.KeyA):
		p.DirX = -1
	case ebiten.IsKeyPressed(ebiten.KeyRight), ebiten.IsKeyPressed(ebiten.KeyD):
		p.DirX = 1
	default:
		moveIntent = false
	}
	// Check jump
	if p.Grounded {
		// Set running state
		if moveIntent {
			st := StateRunning
			if p.State != st {
				p.State, p.StateTick = st, 0
			}
		} else if p.State != StateIdle {
			p.State, p.StateTick = StateIdle, 0
		}
		// Jumping
		switch {
		case inpututil.IsKeyJustPressed(ebiten.KeySpace),
			inpututil.IsKeyJustPressed(ebiten.KeyUp),
			inpututil.IsKeyJustPressed(ebiten.KeyW):
			p.State, p.StateTick = StateJumping, 0
			p.VelocityY = 0
			p.Grounded = false
		}
	}

	// Update position
	if moveIntent {
		p.X += float64(p.DirX) * MovementSpeed
	}
	// Velocity Y
	if !p.Grounded {
		switch p.State {
		case StateJumping:
			p.VelocityY += 1.5
			if p.VelocityY > MaxJumpVelocity {
				p.State, p.StateTick = StateFalling, 0
			}
		case StateFalling:
			p.VelocityY = max(p.VelocityY-1, MaxFallVelocity)
		}
		p.Y += p.VelocityY
	}
	// Ticks
	p.StateTick++
}

func (p *Player) Draw(screen *ebiten.Image) {
	var img *ebiten.Image
	var geom ebiten.GeoM

	frameIndex := (p.StateTick % (FrameDuration * 2)) / FrameDuration
	switch p.State {
	case StateIdle:
		img = ImagePlayerIdle0
		if frameIndex > 0 {
			img = ImagePlayerIdle1
		}
	case StateJumping, StateFalling:
		img = ImagePlayerJump
	case StateRunning:
		img = ImagePlayerRun0
		if frameIndex > 0 {
			img = ImagePlayerRun1
		}
	}
	// Horizontal flip
	if p.DirX == -1 {
		geom.Translate(-PlayerSize/2, 0)
		geom.Scale(-1, 1)
		geom.Translate(PlayerSize/2, 0)
	}
	// Translation
	dx := float64(screen.Bounds().Dx())/2 - PlayerSize/2
	dy := float64(screen.Bounds().Dy()) - PlayerSize
	geom.Translate(dx, dy)
	// Render
	screen.DrawImage(img, &ebiten.DrawImageOptions{
		GeoM: geom,
	})
}
