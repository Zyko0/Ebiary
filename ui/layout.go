package ui

import (
	"image"

	"github.com/Zyko0/Ebiary/ui/internal/graphics"
	"github.com/hajimehoshi/ebiten/v2"
)

type context struct {
	Cursor image.Point

	Z               int
	Hovered         Item
	Unhovered       Item
	DeferredUpdates []func(InputState)
}

// Layout represents a surface (or window) acting as a grid, container
// for child UI items
// A Layout cannot exist without a grid
type Layout struct {
	grid     *Grid
	unit     image.Rectangle
	pipeline *graphics.Pipeline

	offset image.Point
}

func NewLayoutFromGrid(grid *Grid, unit image.Rectangle) *Layout {
	return &Layout{
		grid: grid,
		unit: unit,
		pipeline: graphics.NewPipeline(
			image.Rect(
				0, 0,
				grid.width*unit.Dx(),
				grid.height*unit.Dy(),
			),
		),
	}
}

func NewLayout(columns, rows int, unit image.Rectangle) *Layout {
	return &Layout{
		grid: NewGrid(columns, rows),
		unit: unit,
		pipeline: graphics.NewPipeline(
			image.Rect(
				0, 0,
				columns*unit.Dx(),
				rows*unit.Dy(),
			),
		),
	}
}

func (l *Layout) WithGrid(g *Grid) *Layout {
	l.grid = g
	return l
}

func (l *Layout) Grid() *Grid {
	return l.grid
}

func (l *Layout) Unit() image.Rectangle {
	return l.unit
}

func (l *Layout) SetUnit(unit image.Rectangle) {
	l.unit = unit
}

func (l *Layout) Dimensions() (int, int) {
	return l.grid.width * l.unit.Dx(), l.grid.height * l.unit.Dy()
}

func (l *Layout) SetDimensions(width, height int) {
	w := float64(width) / float64(l.grid.width)
	h := float64(height) / float64(l.grid.height)
	l.unit = image.Rect(0, 0, int(w), int(h))
}

func (l *Layout) Update(offset image.Point, is InputState) {
	l.offset = offset
	width := l.grid.width * l.unit.Dx()
	height := l.grid.height * l.unit.Dy()
	area := image.Rect(
		0, 0,
		width, height,
	)
	c := &context{
		Cursor: is.Cursor().Add(l.offset.Mul(-1)),

		Z: -1,
	}
	l.grid.update(c, area, 0)
	// Update focused items and trigger callbacks
	// TODO: Consider setters for instead of type assertion OR
	// TODO: Consider having ibase instead of items in context
	if c.Hovered != nil {
		c.Hovered.base().state.hoveredTicks++
	}
	if c.Unhovered != nil {
		c.Unhovered.base().state.justUnhovered = true
		c.Unhovered.base().state.hoveredTicks = 0
	}
	// Execute deferred update functions (post hover state processing)
	for _, f := range c.DeferredUpdates {
		f(is)
	}
}

func (l *Layout) Draw(screen *ebiten.Image) {
	l.pipeline.Clear()

	width := l.grid.width * l.unit.Dx()
	height := l.grid.height * l.unit.Dy()
	area := image.Rect(
		0, 0,
		width, height,
	)
	l.grid.addGFX(l.pipeline, area, 0)

	// Screen
	l.pipeline.Draw(screen, l.offset)
}
