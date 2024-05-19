package uiex

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Image is an implementation of Content
type Image struct {
	content

	image *ebiten.Image
	color ebiten.ColorScale
	fill  bool
}

func NewImage(img *ebiten.Image) *Image {
	return &Image{
		content: newContent(),

		image: img,
	}
}

func (i *Image) Draw(dst *ebiten.Image) {
	area := dst.Bounds()

	opts := &ebiten.DrawImageOptions{}
	if i.fill {
		fx := float64(dst.Bounds().Dx()) / float64(i.image.Bounds().Dx())
		fy := float64(dst.Bounds().Dy()) / float64(i.image.Bounds().Dy())
		opts.GeoM.Scale(fx, fy)
		i.lastRegion = area
		i.lastFullRegion = area
	} else {
		x, y := i.alignOffset(
			area, image.Rect(
				0, 0, i.image.Bounds().Dx(), i.image.Bounds().Dy(),
			),
		)
		area = area.Add(image.Pt(int(x), int(y)))
		area = area.Add(i.srcOffset)
		i.lastFullRegion = i.image.Bounds().Add(area.Min)
		area = image.Rect(
			min(max(area.Min.X, i.lastFullRegion.Min.X), area.Max.X),
			min(max(area.Min.Y, i.lastFullRegion.Min.Y), area.Max.Y),
			max(min(area.Max.X, i.lastFullRegion.Max.X), area.Min.X),
			max(min(area.Max.Y, i.lastFullRegion.Max.Y), area.Min.Y),
		)
		i.lastRegion = area
	}

	opts.ColorScale = i.color
	opts.GeoM.Translate(float64(area.Min.X), float64(area.Min.Y))
	dst.DrawImage(i.image, opts)
}

func (i *Image) Image() *ebiten.Image {
	return i.image
}

func (i *Image) SetImage(img *ebiten.Image) {
	i.image = img
}

func (i *Image) SetColorScale(cs ebiten.ColorScale) {
	i.color = cs
}

func (i *Image) SetContainerFilling(fill bool) {
	i.fill = fill
}

// Option

type ImageOption func(*Image)

func (i *Image) WithOptions(opts ...ImageOption) *Image {
	for _, o := range opts {
		o(i)
	}
	return i
}
