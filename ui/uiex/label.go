package uiex

import (
	"image"

	"github.com/Zyko0/Ebiary/ui"
)

type Label struct {
	*block

	text Text
}

func NewLabel(str string) *Label {
	l := &Label{
		text: NewBasicText(str),
	}
	l.block = newBlock(l)
	l.block.SetShape(ui.ShapeNone)
	l.block.SetContent(l.text)
	defaultTheme.apply(l)

	return l
}

func (l *Label) Text() Text {
	return l.text
}

func (l *Label) SetText(txt Text) {
	l.text = txt
	l.block.SetContent(txt)
}

func (l *Label) LastFullRegion() image.Rectangle {
	return l.text.LastFullRegion().Add(l.block.LastFullRegion().Min)
}

func (l *Label) LastRegion() image.Rectangle {
	return l.text.LastRegion().Add(l.block.LastRegion().Min)
}

func (l *Label) SourceOffset() image.Point {
	return l.text.SourceOffset()
}

func (l *Label) SetSourceOffset(offset image.Point) {
	l.text.SetSourceOffset(offset)
}

func (l *Label) SetContent(txt Text) {
	l.SetText(txt)
}

// Options

type LabelOption func(*Label)

func (l *Label) WithOptions(opts ...LabelOption) *Label {
	for _, o := range opts {
		o(l)
	}
	return l
}
