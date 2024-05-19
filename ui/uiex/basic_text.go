package uiex

import (
	"image"
	"image/color"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// BasicText is an implementation of Text which
// comes with limited options.
type BasicText struct {
	textBase

	glyphs []text.Glyph

	color    ebiten.ColorScale
	procText string
}

func NewBasicText(str string) *BasicText {
	bt := &BasicText{
		textBase: newTextBase(),
	}
	bt.SetText(str)

	return bt
}

func (t *BasicText) Draw(dst *ebiten.Image) {
	if t.face == nil {
		panic("uiex: text content must have a font face")
	}

	w, h := text.Measure(t.procText, t.face, t.layout.LineSpacing)
	w, h = math.Round(w+0.5), math.Round(h+0.5)
	area := dst.Bounds()
	area = t.adjustAreaPadding(area)
	inner := image.Rect(0, 0, int(w), int(h)).Add(area.Min)
	ax, ay := t.alignOffset(area, inner)
	area = area.Add(image.Pt(int(ax), int(ay)))
	area = area.Add(t.srcOffset)
	t.lastFullRegion = image.Rect(0, 0, int(w), int(h)).Add(area.Min)
	area = image.Rect(
		min(max(area.Min.X, t.lastFullRegion.Min.X), area.Max.X),
		min(max(area.Min.Y, t.lastFullRegion.Min.Y), area.Max.Y),
		max(min(area.Max.X, t.lastFullRegion.Max.X), area.Min.X),
		max(min(area.Max.Y, t.lastFullRegion.Max.Y), area.Min.Y),
	)
	t.lastRegion = area

	opts := &text.DrawOptions{}
	opts.LayoutOptions = t.layout
	opts.ColorScale = t.color
	opts.GeoM.Translate(float64(area.Min.X), float64(area.Min.Y))
	text.Draw(dst, t.procText, t.face, opts)
}

func (t *BasicText) Glyphs() []text.Glyph {
	if t.updated {
		t.glyphs = text.AppendGlyphs(t.glyphs[:0], t.procText, t.face, &t.layout)
		t.updated = false
	}
	return t.glyphs
}

func (t *BasicText) ProcessedText() string {
	return t.procText
}

func (t *BasicText) SetText(str string) {
	t.textBase.SetText(str)
	// Process rendered text
	spaces := ""
	for i := 0; i < t.tabSpaces; i++ {
		spaces += " "
	}
	t.procText = strings.ReplaceAll(str, "\t", spaces)
}

func (t *BasicText) SetColor(clr color.Color) {
	t.color.Reset()
	t.color.ScaleWithColor(clr)
}

func (t *BasicText) SetColorScale(cs ebiten.ColorScale) {
	t.color = cs
}

// Option

type BasicTextOption func(*BasicText)

func (t *BasicText) WithOptions(opts ...BasicTextOption) *BasicText {
	for _, o := range opts {
		o(t)
	}
	return t
}
