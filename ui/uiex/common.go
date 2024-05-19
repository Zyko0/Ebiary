package uiex

import (
	"image"

	"github.com/Zyko0/Ebiary/ui"
)

// Unexported aliases

type (
	iclass interface {
		Classes() []string
		SetClasses(classes ...string)
		AddClasses(classes ...string)
	}
	_block = ui.Block
	_grid  = ui.Grid
)

type block struct {
	*_block
	addr ui.Item
}

func (b *block) Addr() ui.Item {
	return b.addr
}

func newBlock(addr ui.Item) *block {
	return &block{
		_block: ui.NewBlock(),
		addr:   addr,
	}
}

type grid struct {
	*_grid
	addr ui.Item
}

func (g *grid) Addr() ui.Item {
	return g.addr
}

func newGrid(columns, rows int, addr ui.Item) *grid {
	return &grid{
		_grid: ui.NewGrid(columns, rows),
		addr:  addr,
	}
}

// Classes

type classImpl struct {
	classes []string
}

func (i *classImpl) Classes() []string {
	return i.classes
}

func (i *classImpl) SetClasses(classes ...string) {
	i.classes = append(i.classes[:0], classes...)
}

func (i *classImpl) AddClasses(classes ...string) {
	i.classes = append(i.classes, classes...)
}

// Content

type content struct {
	iclass
	srcOffset      image.Point
	lastFullRegion image.Rectangle
	lastRegion     image.Rectangle
	alignX         ui.AlignMode
	alignY         ui.AlignMode
	paddingLeft    int
	paddingRight   int
	paddingTop     int
	paddingBottom  int
}

func newContent() content {
	return content{
		iclass: &classImpl{},
	}
}

func (c *content) adjustAreaPadding(area image.Rectangle) image.Rectangle {
	area.Min.X += c.paddingLeft
	area.Min.Y += c.paddingTop
	area.Max.X -= c.paddingRight
	area.Max.Y -= c.paddingBottom

	return area
}

func (c *content) alignOffset(area, inner image.Rectangle) (float64, float64) {
	var x, y float64
	switch c.alignX {
	case ui.AlignMin:
		x = 0
	case ui.AlignMax:
		x = float64(area.Max.X - inner.Dx() - inner.Min.X)
	case ui.AlignCenter:
		x = float64(area.Dx())/2 - float64(inner.Dx())/2
	}
	switch c.alignY {
	case ui.AlignMin:
		y = 0
	case ui.AlignMax:
		y = float64(area.Max.Y - inner.Dy() - area.Min.Y)
	case ui.AlignCenter:
		y = float64(area.Dy()/2) - float64(inner.Dy())/2
	}

	return x, y
}

func (c *content) LastFullRegion() image.Rectangle {
	return c.lastFullRegion
}

func (c *content) LastRegion() image.Rectangle {
	return c.lastRegion
}

func (c *content) SourceOffset() image.Point {
	return c.srcOffset
}

func (c *content) SetSourceOffset(offset image.Point) {
	c.srcOffset = offset
}

func (c *content) SetAlignX(x ui.AlignMode) {
	c.alignX = x
}

func (c *content) SetAlignY(y ui.AlignMode) {
	c.alignY = y
}

func (c *content) Padding() (left, right, top, bottom int) {
	return c.paddingLeft, c.paddingRight, c.paddingTop, c.paddingBottom
}

func (c *content) SetPaddingLeft(pixels int) {
	c.paddingLeft = pixels
}

func (c *content) SetPaddingRight(pixels int) {
	c.paddingRight = pixels
}

func (c *content) SetPaddingTop(pixels int) {
	c.paddingTop = pixels
}

func (c *content) SetPaddingBottom(pixels int) {
	c.paddingBottom = pixels
}
