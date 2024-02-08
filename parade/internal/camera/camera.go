package camera

import (
	"math"

	"github.com/go-gl/mathgl/mgl64"
)

// Orthographic

type CameraOrthographic struct {
	position  mgl64.Vec3
	direction mgl64.Vec3

	proj mgl64.Mat4
	view mgl64.Mat4
}

func NewOrthographic(width, height, depth float64) *CameraOrthographic {
	return &CameraOrthographic{
		position:  mgl64.Vec3{0, 0, 0},
		direction: mgl64.Vec3{0, 0, 1},

		proj: mgl64.Ortho(
			-width/2, width/2,
			height/2, -height/2,
			0, depth,
		),
	}
}

func (co *CameraOrthographic) ProjectionMatrix() []float64 {
	return co.proj[:]
}

func (co *CameraOrthographic) ViewMatrix() []float64 {
	return co.view[:]
}

func (co *CameraOrthographic) Position() (float64, float64, float64) {
	return co.position[0], co.position[1], co.position[2]
}

func (co *CameraOrthographic) SetPosition(x, y, z float64) {
	co.position[0] = x
	co.position[1] = y
	co.position[2] = z
}

func (co *CameraOrthographic) Direction() (float64, float64, float64) {
	return co.direction[0], co.direction[1], co.direction[2]
}

func (co *CameraOrthographic) SetDirection(x, y, z float64) {
	co.direction[0] = x
	co.direction[1] = y
	co.direction[2] = z
}

func (co *CameraOrthographic) Update() {
	const halfPi = math.Pi / 2

	co.view = mgl64.LookAt(
		0, 0, 0,
		co.direction[0], co.direction[1], co.direction[2],
		0, -1, 0,
	)
}
