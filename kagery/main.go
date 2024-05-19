package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"runtime"

	"github.com/Zyko0/Ebiary/kagery/assets"
	core "github.com/Zyko0/Ebiary/kagery/core"
	"github.com/Zyko0/Ebiary/kagery/core/clipboard"
	"github.com/Zyko0/Ebiary/kagery/utils"
	"github.com/Zyko0/Ebiary/ui"

	"github.com/Zyko0/Ebiary/ui/opt"
	_ "github.com/Zyko0/Ebiary/ui/theme"
	"github.com/Zyko0/Ebiary/ui/uiex"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	Width, Height = 960, 540

	TPS      = 60
	Columns  = 72
	Rows     = 40
	CellSize = 32

	UIPadding = 16
)

type App struct {
	paused bool
	hidden bool
	aa     bool
	ticks  uint64

	width  int
	height int

	offscreen *ebiten.Image

	fps    *uiex.Label
	layout *ui.Layout
	editor *core.Editor
	images *core.ImageBar

	updated bool
}

func New() *App {
	w, h := ebiten.Monitor().Size()

	app := &App{
		width:  w,
		height: h,

		offscreen: ebiten.NewImage(w, h),
	}

	layout := ui.NewLayout(Columns, Rows, image.Rect(0, 0, CellSize, CellSize))
	layout.SetDimensions(w-UIPadding*2, h-UIPadding*2)
	layout.Grid().SetItemOptions(
		opt.Filling(ui.ColorFillingNone),
		opt.PaddingBottom(1),
	)

	// Top buttons
	layout.Grid().Add(0, 0, 8, 3, core.NewButton(
		"Hide UI", color.RGBA{0, 0, 0, 255}, color.RGBA{32, 32, 32, 255},
	).WithOptions(opt.ButtonText.Options(
		opt.EventAction(ui.EventOptions{
			ui.ReleaseHover: func(i ui.Item) {
				app.hidden = !app.hidden
				if app.hidden {
					i.(*uiex.ButtonText).Text().SetText("Show UI")
				} else {
					i.(*uiex.ButtonText).Text().SetText("Hide UI")
				}
				layout.Grid().ForEach(func(gi ui.Item) {
					if gi == i {
						return
					}
					gi.SetSkipped(app.hidden)
				})
			},
		}),
	)))
	layout.Grid().Add(9, 0, 4, 3, core.NewButton(
		"॥", color.RGBA{0, 0, 0, 255}, color.RGBA{32, 32, 32, 255},
	).WithOptions(opt.ButtonText.Options(
		opt.EventAction(ui.EventOptions{
			ui.ReleaseHover: func(i ui.Item) {
				app.paused = !app.paused
				if app.paused {
					i.(*uiex.ButtonText).Text().SetText(">")
				} else {
					i.(*uiex.ButtonText).Text().SetText("॥")
				}
			},
		}),
	)))
	layout.Grid().Add(14, 0, 8, 3, core.NewButton(
		"AA ×2", color.RGBA{0, 0, 0, 255}, color.RGBA{32, 32, 32, 255},
	).WithOptions(opt.ButtonText.Options(
		opt.EventAction(ui.EventOptions{
			ui.ReleaseHover: func(i ui.Item) {
				app.aa = !app.aa
				if app.aa {
					i.(*uiex.ButtonText).Text().SetText("No AA")
				} else {
					i.(*uiex.ButtonText).Text().SetText("AA ×2")
				}
			},
		}),
	)))
	layout.Grid().Add(23, 0, 8, 3, core.NewButton(
		"00:00.00", color.RGBA{0, 0, 0, 255}, color.RGBA{32, 32, 32, 255},
	).WithOptions(opt.ButtonText.Options(
		ui.WithCustomUpdateFunc(func(b *uiex.ButtonText, is ui.InputState) {
			if b.JustPressed() {
				app.ticks = 0
			}
			mins := app.ticks / TPS / 60
			sec := app.ticks / TPS % 60
			fsec := int(float64(app.ticks%TPS) * (100. / 60.))
			b.Text().SetText(fmt.Sprintf("%02d:%02d.%02d", mins, sec, fsec))
		}),
	)))

	fpsLabel := uiex.NewLabel("0").WithOptions(
		opt.Label.Options(
			opt.Filling(ui.ColorFillingDistance),
			opt.Color(color.Black),
			opt.Alpha(core.Alpha),
			opt.Border(1, color.RGBA{128, 128, 128, 255}),
			opt.Shape(ui.ShapeOctogon),
		),
		opt.Label.Text(
			opt.Text.Color(color.RGBA{192, 255, 192, 255}),
		),
	)
	layout.Grid().Add(32, 0, 8, 3, fpsLabel)

	layout.Grid().Add(Columns-12, 0, 6, 3, core.NewButton(
		"Copy", color.RGBA{255, 96, 0, 255}, color.RGBA{255, 128, 0, 255},
	).WithOptions(opt.ButtonText.Options(
		opt.Border(4, color.RGBA{255, 64, 0, 255}),
		opt.EventAction(ui.EventOptions{
			ui.ReleaseHover: func(_ ui.Item) {
				clipboard.WriteText(app.editor.Text().Text())
			},
		}),
	)))
	layout.Grid().Add(Columns-6, 0, 6, 3, core.NewButton(
		"Paste", color.RGBA{96, 0, 255, 255}, color.RGBA{128, 0, 255, 255},
	).WithOptions(opt.ButtonText.Options(
		opt.Border(4, color.RGBA{64, 0, 255, 255}),
		opt.EventAction(ui.EventOptions{
			ui.ReleaseHover: func(_ ui.Item) {
				app.editor.Text().SetText(clipboard.ReadText())
			},
		}),
	)))

	logger := core.NewLogger()
	layout.Grid().Add(0, Rows-9, Columns-13, 9, logger)

	input := core.NewEditor(logger)
	input.Text().SetText(assets.ShaderSrc)
	layout.Grid().Add(0, 4, Columns-13, Rows-10-4, input)

	images := core.NewImageBar(logger)
	layout.Grid().Add(Columns-12, 4, 12, Rows-4, images)

	ui.SetFocusedItem(input)

	app.fps = fpsLabel
	app.layout = layout
	app.editor = input
	app.images = images

	return app
}

func (g *App) Update() error {
	switch {
	case runtime.GOOS != "js" && ebiten.IsKeyPressed(ebiten.KeyEscape):
		return ebiten.Termination
	case inpututil.IsKeyJustPressed(ebiten.KeyF11):
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	g.layout.Update(image.Point{UIPadding, UIPadding}, ui.GetInputState())
	g.editor.Update()

	if !g.paused {
		g.ticks++
	}

	g.updated = true

	return nil
}

func (g *App) Draw(screen *ebiten.Image) {
	g.fps.Text().SetText(fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()))

	w, h := g.layout.Dimensions()
	subOffscreen := g.offscreen.SubImage(image.Rect(0, 0, w+UIPadding, h+UIPadding)).(*ebiten.Image)
	if g.updated {
		subOffscreen.Clear()
		g.layout.Draw(subOffscreen)
	}

	shader := g.editor.Shader()
	if shader != nil {
		vertices, indices := utils.AppendRectVerticesIndices(nil, nil, &utils.RectOpts{
			DstWidth:  float32(screen.Bounds().Dx()),
			DstHeight: float32(screen.Bounds().Dy()),
			SrcWidth:  float32(g.images.Images()[0].Bounds().Dx()),
			SrcHeight: float32(g.images.Images()[0].Bounds().Dy()),
			R:         1,
			G:         1,
			B:         1,
			A:         1,
		})
		x, y := ebiten.CursorPosition()
		opts := &ebiten.DrawTrianglesShaderOptions{
			Images: g.images.Images(),
			Uniforms: map[string]any{
				"Cursor": []float32{float32(x), float32(y)},
				"Time":   float32(g.ticks/TPS) + float32(g.ticks%TPS)/TPS,
				"Resolution": []float32{
					float32(screen.Bounds().Dx()),
					float32(screen.Bounds().Dy()),
				},
			},
			AntiAlias: g.aa,
		}
		screen.DrawTrianglesShader(vertices, indices, shader, opts)
	}

	opts := &ebiten.DrawImageOptions{}
	opts.ColorScale.ScaleAlpha(1)
	screen.DrawImage(subOffscreen, opts)

	g.updated = false
}

func (g *App) Layout(w, h int) (int, int) {
	g.layout.SetDimensions(w-UIPadding*2, h-UIPadding*2)
	lw, lh := g.layout.Dimensions()
	return lw + UIPadding*2, lh + UIPadding*2
}

func main() {
	ebiten.SetTPS(TPS)
	ebiten.SetVsyncEnabled(true)
	ebiten.SetFullscreen(false)
	ebiten.SetWindowSize(Width, Height)
	ebiten.SetWindowSizeLimits(Width, Height, -1, -1)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	ui.SetFontSource(assets.NotoSource)

	g := New()
	if err := ebiten.RunGameWithOptions(g, &ebiten.RunGameOptions{
		GraphicsLibrary: ebiten.GraphicsLibraryOpenGL,
	}); err != nil {
		log.Fatal("main: ", err)
	}
}
