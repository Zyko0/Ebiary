package theme

import (
	"image/color"

	"github.com/Zyko0/Ebiary/ui"
	"github.com/Zyko0/Ebiary/ui/opt"
	"github.com/Zyko0/Ebiary/ui/uiex"
)

var Default = &uiex.ClassTheme{
	Theme: &uiex.Theme{
		Block: []ui.ItemOption{
			opt.RGB(32, 32, 32),
		},
		Grid:       []ui.ItemOption{},
		Button:     []ui.ItemOption{},
		ButtonText: []ui.ItemOption{},
		Label:      []ui.ItemOption{},
		Bar:        []ui.ItemOption{},
		TextInput: []ui.ItemOption{
			opt.Item(
				opt.TextInput.Caret.Options(
					opt.RGB(0, 0, 128),
				),
				opt.TextInput.Caret.Thickness(3),
			),
		},
		// Decorations
		Scrollbar: []ui.ItemOption{
			opt.RGB(32, 32, 32),
			opt.Filling(ui.ColorFillingNone),
			opt.Border(1, color.Black),
			opt.EventStyle(ui.EventOptions{
				ui.Hover: opt.Item(opt.Scrollbar.Cursor.Options(
					opt.DoEvent(ui.Hover),
				)),
			}),
			opt.Item(
				opt.Scrollbar.Bar.Options(
					opt.Bar.Direction(uiex.DirectionVertical),
					opt.Bar.Thickness(24),
				),
				opt.Scrollbar.Cursor.Options(
					opt.Padding(1),
					opt.EventStyle(ui.EventOptions{
						ui.Default: opt.Multi(
							opt.RGB(32, 32, 32),
							opt.AlphaDecay(5./255, 0.25),
						),
						ui.Hover: opt.Multi(
							opt.AlphaIncr(10./255, 0.75),
						),
						ui.Unhover: opt.Noop(),
						ui.PressHover: opt.Multi(
							opt.RGB(48, 48, 48),
							opt.AlphaIncr(10./255, 0.75),
						),
						ui.Press: opt.DoEvent(ui.PressHover),
					}),
				),
			),
		},
		// Content
		BasicText: []func(*uiex.BasicText){},
		RichText:  []func(*uiex.RichText){},
		Image:     []func(*uiex.Image){},
	},
}

func init() {
	uiex.SetClassTheme(Default)
}
