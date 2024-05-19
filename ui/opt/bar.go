package opt

import (
	"github.com/Zyko0/Ebiary/ui"
	"github.com/Zyko0/Ebiary/ui/uiex"
)

type bar struct{}

var Bar bar

func (bar) Options(opts ...ui.ItemOption) uiex.BarOption {
	return func(b *uiex.Bar) {
		b.SetItemOptions(opts...)
	}
}

func (bar) Thickness(thickness int) uiex.BarOption {
	return func(b *uiex.Bar) {
		b.SetThickness(thickness)
	}
}

func (bar) Length(length int) uiex.BarOption {
	return func(b *uiex.Bar) {
		b.SetLength(length)
	}
}

func (bar) Direction(direction uiex.Direction) uiex.BarOption {
	return func(b *uiex.Bar) {
		b.SetDirection(direction)
	}
}
