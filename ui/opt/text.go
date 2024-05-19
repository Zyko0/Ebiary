package opt

import (
	"bytes"
	"fmt"
	"image/color"
	"io"
	"io/fs"
	"os"

	"github.com/Zyko0/Ebiary/ui"
	"github.com/Zyko0/Ebiary/ui/uiex"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/text/language"
)

type itext[T interface {
	*uiex.BasicText | *uiex.RichText
	uiex.Text

	SetFace(text.Face)
	SetSource(*text.GoTextFaceSource)
	SetSize(float64)
	SetDirection(text.Direction)
	SetLanguage(language.Tag)
	SetScript(language.Script)
	SetLayout(text.LayoutOptions)
	SetPaddingLeft(int)
	SetPaddingRight(int)
	SetPaddingTop(int)
	SetPaddingBottom(int)
}] struct{}

// Alignment

func (itext[T, ]) Align(x, y ui.AlignMode) func(T) {
	return func(t T) {
		t.SetAlignX(x)
		t.SetAlignY(y)
	}
}

func (itext[T]) AlignLeft() func(T) {
	return func(t T) {
		t.SetAlignX(ui.AlignMin)
	}
}

func (itext[T]) AlignRight() func(T) {
	return func(t T) {
		t.SetAlignX(ui.AlignMax)
	}
}

func (itext[T]) AlignCenter() func(T) {
	return func(t T) {
		t.SetAlignX(ui.AlignCenter)
		t.SetAlignY(ui.AlignCenter)
	}
}

func (itext[T]) AlignCenterX() func(T) {
	return func(t T) {
		t.SetAlignX(ui.AlignCenter)
	}
}

func (itext[T]) AlignTop() func(T) {
	return func(t T) {
		t.SetAlignY(ui.AlignMin)
	}
}

func (itext[T]) AlignBottom() func(T) {
	return func(t T) {
		t.SetAlignY(ui.AlignMax)
	}
}

func (itext[T]) AlignCenterY() func(T) {
	return func(t T) {
		t.SetAlignY(ui.AlignCenter)
	}
}

// Padding

func (itext[T]) Padding(pixels int) func(T) {
	return func(t T) {
		t.SetPaddingLeft(pixels)
		t.SetPaddingRight(pixels)
		t.SetPaddingTop(pixels)
		t.SetPaddingBottom(pixels)
	}
}

func (itext[T]) PaddingLeft(pixels int) func(T) {
	return func(t T) {
		t.SetPaddingLeft(pixels)
	}
}

func (itext[T]) PaddingRight(pixels int) func(T) {
	return func(t T) {
		t.SetPaddingRight(pixels)
	}
}

func (itext[T]) PaddingTop(pixels int) func(T) {
	return func(t T) {
		t.SetPaddingTop(pixels)
	}
}

func (itext[T]) PaddingBottom(pixels int) func(T) {
	return func(t T) {
		t.SetPaddingBottom(pixels)
	}
}

// Font

func (itext[T]) Face(face text.Face) func(T) {
	return func(t T) {
		t.SetFace(face)
	}
}

func (itext[T]) Source(source *text.GoTextFaceSource) func(T) {
	return func(t T) {
		t.SetSource(source)
	}
}

func (itext[T]) Size(size float64) func(T) {
	return func(t T) {
		t.SetSize(size)
	}
}

func (itext[T]) Direction(direction text.Direction) func(T) {
	return func(t T) {
		t.SetDirection(direction)
	}
}

func (itext[T]) Language(language language.Tag) func(T) {
	return func(t T) {
		t.SetLanguage(language)
	}
}

func (itext[T]) Script(script language.Script) func(T) {
	return func(t T) {
		t.SetScript(script)
	}
}

// Font parsing

func (itext[T]) FontReader(reader io.Reader) func(T) {
	s, err := text.NewGoTextFaceSource(reader)
	if err != nil {
		panic("uiex: unable to parse text font source: " + err.Error())
	}
	return func(t T) {
		t.SetSource(s)
	}
}

func (i itext[T]) FontPath(path string) func(T) {
	b, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Sprintf("uiex: unable to read font with path %s: %v", path, err))
	}
	return i.FontReader(bytes.NewReader(b))
}

func (i itext[T]) FontFSPath(fsys fs.FS, path string) func(T) {
	f, err := fsys.Open(path)
	if err != nil {
		panic(fmt.Sprintf("uiex: unable to open font at filesystem path %s: %v", path, err))
	}
	return i.FontReader(f)
}

// Layout

func (itext[T]) Layout(layout text.LayoutOptions) func(T) {
	return func(t T) {
		t.SetLayout(layout)
	}
}

func (itext[T]) LineSpacing(pixels float64) func(T) {
	return func(t T) {
		layout := t.Layout()
		layout.LineSpacing = pixels
		t.SetLayout(layout)
	}
}

// Basic text

type basicText struct {
	itext[*uiex.BasicText]
}

var Text basicText

// Classes

func (basicText) Classes(classes ...string) uiex.BasicTextOption {
	return func(t *uiex.BasicText) {
		t.SetClasses(classes...)
	}
}

// Colors

func (basicText) Color(clr color.Color) uiex.BasicTextOption {
	return func(t *uiex.BasicText) {
		t.SetColor(clr)
	}
}

func (basicText) RGBA(r, g, b, a uint8) uiex.BasicTextOption {
	return func(t *uiex.BasicText) {
		t.SetColor(color.RGBA{r, g, b, a})
	}
}

func (basicText) RGB(r, g, b uint8) uiex.BasicTextOption {
	return Text.RGBA(r, g, b, 255)
}

func (basicText) ColorScale(cs ebiten.ColorScale) uiex.BasicTextOption {
	return func(t *uiex.BasicText) {
		t.SetColorScale(cs)
	}
}

// Rich text

type richText struct {
	itext[*uiex.RichText]
}

var RichText richText

// Item text options

type textItem[T interface {
	ui.Item
	Text() uiex.Text
	SetText(uiex.Text)
}] struct{}

func (textItem[T]) SetText(txt uiex.Text) func(T) {
	return func(t T) {
		t.SetText(txt)
	}
}

func (textItem[T]) RichText(opts ...uiex.RichTextOption) func(T) {
	return func(t T) {
		txt := t.Text()
		richtxt, ok := txt.(*uiex.RichText)
		if !ok {
			richtxt = uiex.NewRichText()
			richtxt.SetText(txt.Text())
		}
		for _, o := range opts {
			o(richtxt)
		}
		t.SetText(richtxt)
	}
}

func (textItem[T]) Text(opts ...uiex.BasicTextOption) func(T) {
	return func(t T) {
		txt := t.Text()
		basictxt, ok := txt.(*uiex.BasicText)
		if !ok {
			basictxt = uiex.NewBasicText(txt.Text())
		}
		for _, o := range opts {
			o(basictxt)
		}
		t.SetText(basictxt)
	}
}
