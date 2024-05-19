package uiex

import (
	"github.com/Zyko0/Ebiary/ui"
)

type Bar struct {
	*block
	source ui.Decoration

	visible   bool
	direction Direction
	thickness int
	length    int
}

func NewBar() *Bar {
	b := &Bar{
		visible: true,
	}
	b.block = newBlock(b)
	defaultTheme.apply(b)

	return b
}

func (b *Bar) Visible() bool {
	return b.visible
}

func (b *Bar) SetVisible(visible bool) {
	b.visible = visible
}

func (b *Bar) adjustSize() {
	switch b.direction {
	case DirectionHorizontal:
		b.SetMinHeight(b.thickness)
		b.SetMaxHeight(b.thickness)
		b.SetMinWidth(b.length)
		b.SetMaxWidth(b.length)
	case DirectionVertical:
		b.SetMinHeight(b.length)
		b.SetMaxHeight(b.length)
		b.SetMinWidth(b.thickness)
		b.SetMaxWidth(b.thickness)
	}
}

func (b *Bar) SetThickness(thickness int) {
	b.thickness = thickness
	b.adjustSize()
}

func (b *Bar) SetLength(length int) *Bar {
	b.length = length
	b.adjustSize()
	return b
}

func (b *Bar) SetDirection(direction Direction) *Bar {
	b.direction = direction
	switch b.direction {
	case DirectionHorizontal:
		b.SetAlignX(ui.AlignMin)
	case DirectionVertical:
		b.SetAlignY(ui.AlignMin)
	}
	b.adjustSize()
	return b
}

// Options

type BarOption func(*Bar)

func (b *Bar) WithOptions(opts ...BarOption) *Bar {
	for _, o := range opts {
		o(b)
	}
	return b
}
