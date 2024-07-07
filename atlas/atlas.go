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
	// MinSize is the minimum size of images on the atlas.
	// It is a hint to improve the performance of allocations
	// of new images.
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
				// TODO: make the unmanaged optional
				Unmanaged: false, //true,
			},
		),
		set: packing.NewSet(width, height, setOpts),
	}
}

// Image returns the full atlas image.
func (a *Atlas) Image() *ebiten.Image {
	return a.native
}

// Bounds returns the bounds of the atlas.
func (a *Atlas) Bounds() image.Rectangle {
	return a.native.Bounds()
}

// NewImage allocates a rectangle area on the atlas and
// returns the corresponding image.
// It returns a nil if there was no space to allocate the region
// specified by width, height.
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

// Free is not implemented
func (a *Atlas) Free(img *Image) {
	panic("unimplemented")
	img.Image().Clear()
	a.set.Free(img.bounds)
}
