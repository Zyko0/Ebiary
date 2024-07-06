package atlas

import (
	"image"

	"github.com/Zyko0/Ebiary/atlas/internal/packing"
	"github.com/hajimehoshi/ebiten/v2"
)

type Atlas struct {
	native *ebiten.Image
	set    *packing.Set
}

type NewAtlasOptions struct {
	MinSize image.Point
}

func New(width, height int, opts *NewAtlasOptions) *Atlas {
	var setOpts *packing.NewSetOptions
	if opts != nil {
		setOpts = &packing.NewSetOptions{
			MinSize: opts.MinSize,
		}
	}

	return &Atlas{
		native: ebiten.NewImageWithOptions(
			image.Rect(0, 0, width, height),
			&ebiten.NewImageOptions{
				Unmanaged: true,
			},
		),
		set: packing.NewSet(width, height, setOpts),
	}
}

func (a *Atlas) Image() *ebiten.Image {
	return a.native
}

func (a *Atlas) Bounds() image.Rectangle {
	return a.native.Bounds()
}

func (a *Atlas) NewImage(width, height int) *Image {
	r := image.Rect(0, 0, width, height)
	img := &Image{
		atlas:  a,
		bounds: &r,
	}
	if !a.set.Insert(img.bounds) {
		return nil
	}

	return img
}

func (a *Atlas) Free(img *Image) {
	img.Image().Clear()
	a.set.Free(img.bounds)
}
