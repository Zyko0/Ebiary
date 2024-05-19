package core

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"io"
	"io/fs"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/Zyko0/Ebiary/kagery/assets"
	"github.com/Zyko0/Ebiary/kagery/core/clipboard"
	"github.com/Zyko0/Ebiary/ui"
	"github.com/Zyko0/Ebiary/ui/opt"
	"github.com/Zyko0/Ebiary/ui/uiex"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Editor struct {
	*uiex.TextInput

	lexer  *Lexer
	buffer *buffer

	linesText *uiex.RichText
	lines     *uiex.Bar
	logger    *Logger

	lastText string
	shader   *ebiten.Shader
	errLine  int
	err      error
}

func NewEditor(logger *Logger) *Editor {
	lines := uiex.NewBar().WithOptions(
		opt.Bar.Direction(uiex.DirectionVertical),
		opt.Bar.Thickness(48),
		opt.Bar.Options(
			opt.RGB(16, 16, 16),
			opt.Alpha(Alpha),
			opt.AlignLeft(),
			opt.Rounding(Rounding),
		),
	)
	lines.SetFocusHandling(false)
	linesText := uiex.NewRichText().WithOptions(
		opt.RichText.Size(16),
		opt.RichText.Align(ui.AlignOffset, ui.AlignOffset),
		opt.RichText.PaddingLeft(32),
		opt.RichText.Layout(text.LayoutOptions{
			LineSpacing:    18,
			PrimaryAlign:   text.AlignEnd,
			SecondaryAlign: text.AlignStart,
		}),
	)
	lines.SetContent(linesText)

	return &Editor{
		TextInput: uiex.NewTextInput().WithOptions(
			opt.TextInput.Caret.Options(
				opt.RGB(255, 255, 255),
			),
			opt.TextInput.RichText(
				opt.RichText.Source(assets.SourceCodeProSource),
				opt.RichText.Size(16),
				opt.RichText.LineSpacing(18),
				opt.RichText.PaddingLeft(48),
				opt.RichText.PaddingRight(24),
				opt.RichText.PaddingBottom(24),
			),
			opt.TextInput.Options(
				opt.RGB(0, 0, 0),
				opt.Alpha(Alpha),
				opt.Rounding(Rounding),
				// Hack: margin -2 for the border not to be drawn over
				// by lines and scrollbars
				opt.Margin(-2),
				opt.Border(1, clrBorder),
				opt.Decorations(
					lines,
					NewScrollbar(uiex.DirectionHorizontal).WithOptions(
						// Some space for the second scrollbar overlap
						opt.Scrollbar.Options(
							opt.PaddingRight(24),
							opt.PaddingLeft(48),
						),
					),
					NewScrollbar(uiex.DirectionVertical).WithOptions(
						// Some space for the second scrollbar overlap
						opt.Scrollbar.Options(opt.PaddingBottom(24)),
					),
				),
			),
		),

		lexer:  newLexer(),
		buffer: newBuffer(),

		linesText: linesText,
		lines:     lines,
		logger:    logger,
	}
}

func (e *Editor) handleFileDrop() {
	if !e.Hovered() {
		return
	}

	files := ebiten.DroppedFiles()
	if files == nil {
		return
	}
	var txt string
	var found bool
	var txtPath string
	err := fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || found {
			return nil
		}

		txtPath = path
		if err != nil {
			return err
		}
		f, err := files.Open(path)
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		const maxSize = 1024 * 1024
		if info.Size() > maxSize {
			return errors.New("file is too big")
		}
		b, err := io.ReadAll(f)
		if err != nil {
			return err
		}
		if !utf8.Valid(b) {
			return errors.New("text is not valid UTF-8")
		}
		txt = string(b)
		found = true
		return nil
	})
	if err != nil {
		e.logger.RegisterError(&shaderFileError{
			errorBase: &errorBase{
				msg: err.Error() + ": " + txtPath,
			},
		})
		return
	}
	if found {
		e.Text().SetText(txt)
	}
}

func (e *Editor) newShader(txt string) *ebiten.Shader {
	e.err = nil
	e.errLine = -1
	defer func() {
		if r := recover(); r != nil {
			e.err = &panicError{
				&errorBase{
					msg: fmt.Sprintf("%v", r),
				},
			}
		}
	}()
	s, err := ebiten.NewShader([]byte(txt))
	if err != nil {
		errPos := regex.FindAllStringSubmatch(err.Error(), 3)
		if len(errPos) > 0 && len(errPos[0]) > 0 {
			lineNum, _ := strconv.ParseInt(errPos[0][1], 10, 64)
			lineChar, _ := strconv.ParseInt(errPos[0][2], 10, 64)
			e.err = &shaderError{
				errorBase: &errorBase{err.Error()},
				line:      int(lineNum),
				char:      int(lineChar),
			}
			e.errLine = int(lineNum)
		}
	}

	return s
}

var regex = regexp.MustCompile(`(\d+):(\d+):.+`)

func (e *Editor) Update() {
	// Ctrl commands
	var skipNewBuffer bool
	keyPresses := inpututil.AppendJustPressedKeys(nil)
	if ebiten.IsKeyPressed(ebiten.KeyControl) {
		switch {
		// Paste
		case inpututil.IsKeyJustPressed(ebiten.KeyV):
			e.InsertAtCaret(clipboard.ReadText())
		// Start of line
		case inpututil.IsKeyJustPressed(ebiten.KeyQ):
			pos := e.CaretPosition()
			idx := strings.LastIndex(e.Text().Text()[:pos], "\n") + 1
			e.SetCaretPosition(idx)
		// End of line
		case inpututil.IsKeyJustPressed(ebiten.KeyE):
			pos := e.CaretPosition()
			idx := strings.Index(e.Text().Text()[pos:], "\n")
			if idx == -1 {
				idx = len(e.Text().Text())
			} else {
				idx += pos
			}
			e.SetCaretPosition(idx)
		// Cancel/Redo buffer
		case !e.buffer.Empty():
			for _, k := range keyPresses {
				switch ebiten.KeyName(k) {
				// Cancel
				case "z":
					e.Text().SetText(e.buffer.Cancel())
					skipNewBuffer = true
				// Redo
				case "y":
					e.Text().SetText(e.buffer.Redo())
					skipNewBuffer = true
				}
				if skipNewBuffer {
					break
				}
			}
		}
	}
	// Handle dropped file
	e.handleFileDrop()

	txt := e.Text().Text()

	// On text changes
	if txt != e.lastText {
		// Update syntaxic colorization
		err := e.lexer.Format(e.Text().(*uiex.RichText))
		_ = err // TODO: log or not
		// Recompile shader
		s := e.newShader(txt)
		if e.err == nil {
			if e.shader != nil {
				e.shader.Deallocate()
			}
			e.shader = s
			e.logger.RegisterError(nil)
			e.SetBorderColor(clrBorder)
		}
		// Cancel / Redo buffers
		if !skipNewBuffer {
			e.buffer.New(txt)
		}
		// Update logger
		if e.err != nil {
			e.logger.RegisterError(e.err)
			e.SetBorderColor(clrBorderError)
		}
		// Lines count
		e.linesText.Reset()
		e.linesText.SetText("")
		e.linesText.PushColorFg(color.RGBA{128, 128, 128, 255})
		for i := 0; i < strings.Count(txt, "\n")+1; i++ {
			line := strconv.FormatInt(int64(i+1), 10) + "\n"
			if i+1 == e.errLine {
				e.linesText.PushColorFg(color.RGBA{255, 0, 0, 255})
				e.linesText.PushBold()
				e.linesText.Append(line)
				e.linesText.Pop()
				e.linesText.Pop()
				continue
			}
			e.linesText.Append(line)
		}
		e.linesText.Pop()
	}
	e.linesText.SetSourceOffset(image.Pt(0, e.SourceOffset().Y))

	e.lastText = txt
}

func (e *Editor) Shader() *ebiten.Shader {
	return e.shader
}
