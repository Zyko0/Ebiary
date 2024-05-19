package uiex

import (
	"image"

	"github.com/Zyko0/Ebiary/ui"
)

type Amovible interface {
	ui.Item
	SourceOffset() image.Point
	SetSourceOffset(offset image.Point)
	LastFullRegion() image.Rectangle
	LastRegion() image.Rectangle
}

type Scrollable interface {
	Amovible
	LastCursor() image.Point
}

type Draggable interface {
}
