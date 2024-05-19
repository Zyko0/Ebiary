package opt

import (
	"github.com/Zyko0/Ebiary/ui"
	"github.com/Zyko0/Ebiary/ui/uiex"
)

type (
	scrollbarCursorBar struct{}

	scrollbarCursor struct {
		Bar scrollbarCursorBar
	}

	scrollbarBar struct{}

	scrollbar struct {
		Bar    scrollbarBar
		Cursor scrollbarCursor
	}
)

var Scrollbar scrollbar

func (scrollbar) Options(opts ...ui.ItemOption) uiex.ScrollbarOption {
	return func(sb *uiex.Scrollbar) {
		sb.SetItemOptions(opts...)
	}
}

func (scrollbar) Direction(direction uiex.Direction) uiex.ScrollbarOption {
	return func(sb *uiex.Scrollbar) {
		sb.SetDirection(direction)
	}
}

func (scrollbarBar) Options(opts ...uiex.BarOption) uiex.ScrollbarOption {
	return func(sb *uiex.Scrollbar) {
		sb.WithBarOptions(opts...)
	}
}

func (scrollbarCursor) Options(opts ...ui.ItemOption) uiex.ScrollbarOption {
	return func(sb *uiex.Scrollbar) {
		sb.WithCursorOptions(opts...)
	}
}

func (scrollbarCursorBar) Options(opts ...uiex.BarOption) uiex.ScrollbarOption {
	return func(sb *uiex.Scrollbar) {
		sb.WithCursorBarOptions(opts...)
	}
}
