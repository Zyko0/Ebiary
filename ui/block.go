package ui

import (
	"image"

	"github.com/Zyko0/Ebiary/ui/internal/graphics"
)

type Block struct {
	*itemImpl

	content Content
}

func NewBlock() *Block {
	b := &Block{}
	b.itemImpl = newItem(b)

	return b
}

func (b *Block) update(c *context, area image.Rectangle, z int) {
	if b.Skipped() {
		return
	}

	clamped, inner := b.adjustInnerArea(area)
	b.itemImpl.geom.lastFullRegion = inner
	b.itemImpl.geom.lastRegion = clamped
	b.itemImpl.geom.lastCursor = c.Cursor
	b.itemImpl.update(c, b, clamped, z)
	// Decorations
	for _, d := range b.itemImpl.decorations {
		if d.Visible() {
			d.update(c, clamped, z)
		}
	}
}

func (b *Block) addGFX(pp *graphics.Pipeline, area image.Rectangle, z int) {
	if b.Skipped() {
		return
	}

	inner, _ := b.adjustInnerArea(area)
	if !inner.In(area) {
		return
	}
	b.itemImpl.addGFX(pp, inner, z)
	// Content
	if b.content != nil {
		pp.AddContent(&graphics.ContentPrimitive{
			Z:       z,
			Clip:    inner,
			Content: b.content,
		})
	}
	// Decorations
	for _, d := range b.itemImpl.decorations {
		if d.Visible() {
			// Note: mark them as z+1 so that they're drawn on top
			// of the block's content
			d.addGFX(pp, inner, z+1)
		}
	}
}

func (b *Block) BlockBase() *Block {
	return b
}

func (b *Block) Content() Content {
	return b.content
}

func (b *Block) SetContent(c Content) {
	b.content = c
}

// Options

type BlockOption = func(*Block)

func (b *Block) WithOptions(opts ...BlockOption) *Block {
	for _, o := range opts {
		o(b)
	}
	return b
}

func WithBlockContent(c Content) BlockOption {
	return func(b *Block) {
		b.SetContent(c)
	}
}
