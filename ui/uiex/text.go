package uiex

import (
	"image"

	"github.com/Zyko0/Ebiary/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/text/language"
)

type Text interface {
	ui.Content

	Face() text.Face
	Layout() text.LayoutOptions
	TabSpaces() int
	Glyphs() []text.Glyph

	LastFullRegion() image.Rectangle
	LastRegion() image.Rectangle
	SourceOffset() image.Point
	SetSourceOffset(offset image.Point)
	SetAlignX(x ui.AlignMode)
	SetAlignY(y ui.AlignMode)
	Padding() (left, right, top, bottom int)

	Text() string
	SetText(str string)
	ProcessedText() string
	Draw(dst *ebiten.Image)
}

type textBase struct {
	content

	tabSpaces int
	source    *text.GoTextFaceSource
	goface    *text.GoTextFace
	face      text.Face
	layout    text.LayoutOptions

	updated bool
	text    string
}

func newTextBase() textBase {
	source := ui.GetFontSource()
	goface := &text.GoTextFace{
		Source: source,
		Size:   ui.GetFontSize(),
	}

	return textBase{
		content: newContent(),

		tabSpaces: 4,
		source:    source,
		goface:    goface,
		face:      goface,
		layout: text.LayoutOptions{
			LineSpacing: goface.Size * 1.5,
		},

		updated: true,
	}
}

func (tb *textBase) initGoFace() {
	if tb.face != nil {
		if goface, ok := tb.face.(*text.GoTextFace); ok {
			tb.goface = goface
		}
	}
	if tb.goface == nil {
		if tb.source == nil {
			panic("uiex: text font source is not set")
		}
		tb.goface = &text.GoTextFace{
			Source: tb.source,
		}
		tb.face = tb.goface
	}
	if tb.goface.Size == 0 {
		tb.goface.Size = ui.GetFontSize()
	}
	if tb.layout.LineSpacing == 0 {
		tb.layout.LineSpacing = tb.goface.Size * 1.5
	}
}

func (tb *textBase) Face() text.Face {
	return tb.face
}

func (tb *textBase) SetTabSpaces(count int) {
	if count <= 0 {
		panic("uiex: text tab spaces count must be > 0")
	}
	tb.tabSpaces = count
	tb.updated = true
}

func (tb *textBase) TabSpaces() int {
	return tb.tabSpaces
}

func (tb *textBase) SetFace(face text.Face) {
	tb.face = face
	tb.updated = true
}

func (tb *textBase) SetSource(source *text.GoTextFaceSource) {
	tb.source = source
	if tb.goface != nil {
		tb.goface.Source = source
	}
	tb.initGoFace()
	tb.updated = true
}

func (tb *textBase) SetSize(size float64) {
	tb.initGoFace()
	tb.goface.Size = size
	if tb.layout.LineSpacing == 0 {
		tb.layout.LineSpacing = size * 1.5
	}
	tb.updated = true
}

func (tb *textBase) SetDirection(direction text.Direction) {
	tb.initGoFace()
	tb.goface.Direction = direction
	tb.updated = true
}

func (tb *textBase) SetLanguage(language language.Tag) {
	tb.initGoFace()
	tb.goface.Language = language
	tb.updated = true
}

func (tb *textBase) SetScript(script language.Script) {
	tb.initGoFace()
	tb.goface.Script = script
	tb.updated = true
}

func (tb *textBase) Layout() text.LayoutOptions {
	return tb.layout
}

func (tb *textBase) SetLayout(opts text.LayoutOptions) {
	tb.layout = opts
	tb.updated = true
}

func (tb *textBase) Text() string {
	return tb.text
}

func (tb *textBase) SetText(str string) {
	// Re-cache glyphs if text changed
	if str != tb.text {
		tb.updated = true
	}
	tb.text = str
}
