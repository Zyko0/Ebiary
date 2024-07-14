package atlas

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Image struct {
	atlas *Atlas

	bounds *image.Rectangle
}

func (i *Image) Atlas() *Atlas {
	return i.atlas
}

func (i *Image) Bounds() image.Rectangle {
	return *i.bounds
}

func (i *Image) Image() *ebiten.Image {
	return i.atlas.native.SubImage(*i.bounds).(*ebiten.Image)
}

func (i *Image) SubImage(bounds image.Rectangle) *Image {
	bounds = bounds.Add(i.bounds.Min)

	return &Image{
		atlas: i.atlas,

		bounds: &bounds,
	}
}
