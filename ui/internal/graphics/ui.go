package graphics

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// Item

type Shape byte

const (
	ShapeBox Shape = iota
	ShapeEllipse
	ShapeRhombus
	ShapeOctogon
	ShapeNone
)

type ColorFilling byte

const (
	ColorFillingVertical ColorFilling = iota
	ColorFillingDistance
	ColorFillingNone
)

type ItemPrimitive struct {
	Z int
	// First vec4
	Shape          Shape
	ColorMin       color.Color
	ColorMax       color.Color
	ColorMinFactor float32
	// Second vec4
	ColorFilling ColorFilling
	Rounding     float32
	BorderColor  color.Color
	BorderWidth  float32
	// CPU
	MarginLeft   float32
	MarginRight  float32
	MarginTop    float32
	MarginBottom float32

	// Extras
	ColorAlpha     float32
}

// Content

type Content interface {
	Draw(dst *ebiten.Image)
}

type ContentPrimitive struct {
	Z       int
	Clip    image.Rectangle
	Content Content
}
