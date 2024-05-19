package core

import (
	"image/color"

	"github.com/Zyko0/Ebiary/ui"
	"github.com/Zyko0/Ebiary/ui/opt"
	"github.com/Zyko0/Ebiary/ui/uiex"
)

const (
	Rounding = 15.
	Alpha    = 0.75
)

var (
	clrBorder      = color.RGBA{200, 200, 200, 255}
	clrBorderError = color.RGBA{255, 0, 0, 255}
)

func NewButton(txt string, clr, hover color.RGBA) *uiex.ButtonText {
	return uiex.NewButtonText(txt).WithOptions(
		opt.ButtonText.Options(
			opt.Rounding(Rounding),
			opt.Color(clr),
			opt.Alpha(Alpha),
			opt.Border(1, clrBorder),
			opt.EventStyle(ui.EventOptions{
				ui.Default: opt.Color(clr),
				ui.Hover:   opt.Color(hover),
			}),
		),
	)
}

func NewScrollbar(direction uiex.Direction) *uiex.Scrollbar {
	return uiex.NewScrollbar().WithOptions(
		opt.Scrollbar.Direction(direction),
		opt.Scrollbar.Options(
			opt.RGB(0, 0, 0),
			opt.Border(1, color.Black),
			opt.Rounding(Rounding),
			opt.EventStyle(ui.EventOptions{
				ui.Hover: opt.Item(opt.Scrollbar.Cursor.Options(
					opt.DoEvent(ui.Hover),
				)),
			}),
		),
		opt.Scrollbar.Cursor.Options(
			opt.Padding(1),
			opt.Rounding(Rounding),
			opt.EventStyle(ui.EventOptions{
				ui.Default: opt.Multi(
					opt.RGB(128, 128, 128),
					opt.AlphaDecay(5./255, 0.25),
				),
				ui.Hover: opt.Multi(
					opt.AlphaIncr(10./255, 0.75),
				),
				ui.Unhover: opt.Noop(),
				ui.PressHover: opt.Multi(
					opt.RGB(192, 192, 192),
					opt.AlphaIncr(10./255, 0.75),
				),
				ui.Press: opt.DoEvent(ui.PressHover),
			}),
		),
	)
}
