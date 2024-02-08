package parade

import "github.com/hajimehoshi/ebiten/v2"

type Layer struct {
	Z        float64
	Depth    float64
	Diffuse  *ebiten.Image
	Height   *ebiten.Image
	Normal   *ebiten.Image
	Specular *ebiten.Image

	// BoxMapped defines whether boxmapping technique should be used to
	// color the volume's texture
	BoxMapped bool
}
