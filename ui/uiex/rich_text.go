package uiex

import (
	"image"
	"image/color"
	"math"
	"sort"
	"strings"

	"github.com/Zyko0/Ebiary/ui/uiex/internal/uitext"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// RichText is an implementation of Text which allows
// various text effects.
type RichText struct {
	textBase

	glyphs []text.Glyph

	stack    *uitext.Stack
	procText string
}

func NewRichText() *RichText {
	rt := &RichText{
		textBase: newTextBase(),

		stack: uitext.NewEffectStack(),
	}
	rt.SetText("")

	return rt
}

func (t *RichText) Draw(dst *ebiten.Image) {
	if t.face == nil {
		panic("uiex: rich text content must have a font face")
	}

	glyphs := t.Glyphs()
	w, h := text.Measure(t.procText, t.face, t.layout.LineSpacing)
	w, h = math.Round(w+0.5), math.Round(h+0.5)
	area := dst.Bounds()
	ori := area
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

	metrics := t.face.Metrics()

	boxIndices := []uint16{0, 1, 2, 1, 2, 3}
	ir := 0
	y := 0
	lh := float32(metrics.HAscent + metrics.HDescent)
	lastGlyphX := float32(0.)
	opts := &ebiten.DrawTrianglesShaderOptions{
		Uniforms: map[string]any{
			"LineHeight": lh,
		},
	}
	geom := ebiten.GeoM{}
	geom.Skew(0.1396, 0)
	t.stack.ResetIndex()
	effects := append([]*uitext.TextEffect{}, t.stack.Effects()...)
	sort.SliceStable(effects, func(i, j int) bool {
		return effects[i].Start < effects[j].Start
	})
	_ = effects
	for _, g := range glyphs {
		for ir < len(t.procText) && t.procText[ir] == '\n' {
			lastGlyphX = 0
			y++
			ir++
		}

		end := ir + (g.EndIndexInBytes - g.StartIndexInBytes)
		effect := t.stack.Effect(ir)
		rw := float32(text.Advance(string(t.procText[ir:end]), t.face))
		gx := float32(g.X) + float32(area.Min.X)
		gy := float32(g.Y) + float32(area.Min.Y)
		xoff := lastGlyphX - float32(g.X)
		if lastGlyphX == 0 {
			// Note: area.Min is the starting X
			xoff = float32(-g.X)
		}
		if effect.Mask&uitext.TextMaskItalic > 0 {
			x, _ := geom.Apply(g.X, g.Y)
			if g.X-x < 0 {
				// TODO: make the background as a separate pass instead of this hack
				rw -= float32(g.X-x) * 5
			}
		}
		xoff = min(xoff, 0)
		if effect.Mask&uitext.TextMaskBold > 0 {
			// Compensate faux bold src offset from shader
			xoff -= lh * 0.025
		}
		yoff := float32(y) * float32(t.layout.LineSpacing)
		yoff = yoff - float32(g.Y)
		// Skip oob glyphs
		if gx+xoff > float32(ori.Max.X) || gy+yoff > float32(ori.Max.Y) ||
			gx+rw < float32(ori.Min.X) || gy+yoff+lh < float32(ori.Min.Y) {
			ir = end
			continue
		}
		mask := float32(effect.Mask)
		fg := uitext.ColorAsFloat32RGB(effect.ColorFg)
		bg := uitext.ColorAsFloat32RGB(effect.ColorBg)
		vertices := []ebiten.Vertex{
			{
				DstX:   gx + xoff,
				DstY:   gy + yoff,
				SrcX:   xoff,
				SrcY:   yoff,
				ColorR: fg,
				ColorG: bg,
				ColorB: 1,
				ColorA: mask,
			},
			{
				DstX:   gx + rw,
				DstY:   gy + yoff,
				SrcX:   rw,
				SrcY:   yoff,
				ColorR: fg,
				ColorG: bg,
				ColorB: 1,
				ColorA: mask,
			},
			{
				DstX:   gx + xoff,
				DstY:   gy + yoff + lh,
				SrcX:   xoff,
				SrcY:   yoff + lh,
				ColorR: fg,
				ColorG: bg,
				ColorB: 1,
				ColorA: mask,
			},
			{
				DstX:   gx + rw,
				DstY:   gy + yoff + lh,
				SrcX:   rw,
				SrcY:   yoff + lh,
				ColorR: fg,
				ColorG: bg,
				ColorB: 1,
				ColorA: mask,
			},
		}
		opts.Images[0] = g.Image
		_, _ = vertices, boxIndices
		dst.DrawTrianglesShader(vertices, boxIndices, uitext.Shader(), opts)
		lastGlyphX = gx + xoff + rw
		ir = end
	}
}

func (t *RichText) Glyphs() []text.Glyph {
	if t.updated {
		t.glyphs = text.AppendGlyphs(t.glyphs[:0], t.procText, t.face, &t.layout)
		t.updated = false
	}
	return t.glyphs
}

func (t *RichText) ProcessedText() string {
	return t.procText
}

func (t *RichText) SetText(str string) {
	t.textBase.SetText(str)
	// Process rendered text
	spaces := ""
	for i := 0; i < t.tabSpaces; i++ {
		spaces += " "
	}
	t.procText = strings.ReplaceAll(str, "\t", spaces)
}

func (t *RichText) Append(str string) {
	t.SetText(t.textBase.text + str)
}

func (t *RichText) PushBold() {
	t.stack.PushMask(len(t.procText), uitext.TextMaskBold)
}

func (t *RichText) PushItalic() {
	t.stack.PushMask(len(t.procText), uitext.TextMaskItalic)
}

func (t *RichText) PushUnderline() {
	t.stack.PushMask(len(t.procText), uitext.TextMaskUnderline)
}

func (t *RichText) PushStrikethrough() {
	t.stack.PushMask(len(t.procText), uitext.TextMaskStrikethrough)
}

func (t *RichText) PushColorFg(clr color.Color) {
	t.stack.PushFg(len(t.procText), clr)
}

func (t *RichText) PushColorBg(clr color.Color) {
	t.stack.PushBg(len(t.procText), clr)
}

func (t *RichText) Pop() {
	t.stack.Pop(len(t.procText))
}

func (t *RichText) Reset() {
	t.stack.Reset()
}

// Option

type RichTextOption func(*RichText)

func (t *RichText) WithOptions(opts ...RichTextOption) *RichText {
	for _, o := range opts {
		o(t)
	}
	return t
}