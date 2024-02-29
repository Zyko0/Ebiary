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
			-10000, 10000,
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

var Angle = 0.

func rotateX(y, z, a float64) (float64, float64) {
	cosa := math.Cos(a)
	sina := math.Sin(a)
	yy := y*cosa - z*sina
	zz := y*sina + z*cosa
	y = yy
	z = zz
	return y, z
}

func rotateY(x, z, a float64) (float64, float64) {
	cosa := math.Cos(a)
	sina := math.Sin(a)
	xx := x*cosa + z*sina
	zz := z*cosa - x*sina
	x = xx
	z = zz
	return x, z
}

func rotateZ(x, y, a float64) (float64, float64) {
	cosa := math.Cos(a)
	sina := math.Sin(a)
	xx := x*cosa - y*sina
	yy := x*sina + y*cosa
	x = xx
	y = yy
	return x, y
}

func (co *CameraOrthographic) Update() {
	const halfPi = math.Pi / 2

	Angle += 0.015
	Angle = math.Mod(Angle, math.Pi*2)
	a := Angle - math.Pi
	tx, ty, tz := co.Position()
	/*x := -1.
	y := -1.
	x = -tx
	y = -ty

	x, y = rotateZ(x, y, a)
	_, _ = x, y*/
	cosa := math.Cos(a)
	sina := math.Sin(a)

	x := -1.
	y := -1.
	nx := x*cosa - y*sina
	ny := y*cosa + x*sina
	x = nx + 1
	y = ny + 1
	_, _ = x, y
	_, _, _ = tx, ty, tz
	pos := mgl64.Vec3{}
	dir := mgl64.Vec3{x, y, 1}.Normalize()
	dir[2] = 1
	co.view = mgl64.LookAtV(
		pos,          //co.position[0], co.position[1], co.position[2], //0, 0, 0,
		pos.Add(dir), //co.direction[0], co.direction[1], co.direction[2],
		mgl64.Vec3{0, -1, 0},
	)
}
