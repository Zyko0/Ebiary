package core

import (
	"image/color"
	"time"

	"github.com/Zyko0/Ebiary/ui"
	"github.com/Zyko0/Ebiary/ui/opt"
	"github.com/Zyko0/Ebiary/ui/uiex"
)

type Logger struct {
	*uiex.Label
}

func NewLogger() *Logger {
	logger := uiex.NewLabel("").WithOptions(
		opt.Label.RichText(
			opt.RichText.Size(18),
			opt.RichText.LineSpacing(20),
			opt.RichText.PaddingRight(24),
			opt.RichText.PaddingBottom(24),
			opt.RichText.PaddingLeft(12),
			opt.RichText.PaddingTop(12),
			opt.RichText.Align(ui.AlignMin, ui.AlignMin),
		),
		opt.Label.Options(
			opt.Shape(ui.ShapeBox),
			opt.RGB(0, 0, 0),
			opt.Alpha(Alpha),
			opt.Rounding(Rounding),
			// Hack: margin -2 for the border not to be drawn over
			// by scrollbars
			opt.Margin(-2),
			opt.Border(1, clrBorder),
			opt.Decorations(
				NewScrollbar(uiex.DirectionHorizontal).WithOptions(
					opt.Scrollbar.Options(opt.PaddingRight(24)),
				),
				NewScrollbar(uiex.DirectionVertical).WithOptions(
					opt.Scrollbar.Options(opt.PaddingBottom(24)),
				),
			),
		),
	)

	return &Logger{
		Label: logger,
	}
}

func (l *Logger) RegisterError(err error) {
	now := "[" + time.Now().Format("15:04:05") + "]"

	rt := l.Text().(*uiex.RichText)
	rt.SetText("")
	rt.Reset()

	if err == nil {
		rt.PushColorFg(color.RGBA{0, 192, 0, 255})
		rt.Append(now + " Compilation successful.\n")
		rt.Pop()
		return
	}

	switch errImpl := err.(type) {
	case *panicError:
		rt.PushColorFg(color.RGBA{255, 128, 0, 255})
		rt.Append(now + " ")
		rt.PushBold()
		rt.Append("Panic! ")
		rt.Pop()
		rt.Append("Please open an issue on https://github.com/hajimehoshi/ebiten if it's not fixed upstream!\n")
		rt.Pop()
		rt.PushColorFg(color.RGBA{255, 0, 0, 255})
		rt.Append(now + " " + errImpl.msg + "\n")
		rt.Pop()
	default:
		//case *shaderError:
		//case *imageError:
		//case *shaderFileError:
		rt.PushColorFg(color.RGBA{192, 0, 0, 255})
		rt.Append(now + " Error: " + err.Error() + "\n")
		rt.Pop()
	}
}
