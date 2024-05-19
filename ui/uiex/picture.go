package uiex

import (
	"image"

	"github.com/Zyko0/Ebiary/ui"
	"github.com/hajimehoshi/ebiten/v2"
)

type Picture struct {
	*block

	img *Image
}

func NewPicture(img *ebiten.Image) *Picture {
	p := &Picture{}
	p.block = newBlock(p)
	p.block.SetColorFilling(ui.ColorFillingNone)
	p.img = NewImage(img)
	p.block.SetContent(p.img)
	defaultTheme.apply(p)

	return p
}

func (p *Picture) Image() *Image {
	return p.img
}

func (p *Picture) LastFullRegion() image.Rectangle {
	return p.img.LastFullRegion()
}

func (p *Picture) LastRegion() image.Rectangle {
	return p.img.LastRegion()
}

func (p *Picture) SourceOffset() image.Point {
	return p.img.SourceOffset()
}

func (p *Picture) SetSourceOffset(offset image.Point) {
	p.img.SetSourceOffset(offset)
}

func (p *Picture) Content() *Image {
	return p.img
}

// Options

type PictureOption func(*Picture)

func (p *Picture) WithOptions(opts ...PictureOption) *Picture {
	for _, o := range opts {
		o(p)
	}
	return p
}
