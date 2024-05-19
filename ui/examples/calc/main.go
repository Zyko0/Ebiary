package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"os"

	"github.com/Zyko0/Ebiary/ui"
	assets "github.com/Zyko0/Ebiary/ui/examples/_assets"
	"github.com/Zyko0/Ebiary/ui/examples/calc/app"
	"github.com/Zyko0/Ebiary/ui/opt"
	"github.com/Zyko0/Ebiary/ui/uiex"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	TPS          = 60
	Columns      = 16
	Rows         = 25
	CellWidth    = 20
	CellHeight   = 20
	ScreenWidth  = Columns * CellWidth
	ScreenHeight = Rows * CellHeight
)

var (
	ebitenImg = ebiten.NewImage(18, 18)

	appTheme = &uiex.ClassTheme{
		Theme: &uiex.Theme{
			BasicText: uiex.StyleContent[*uiex.BasicText]{
				opt.Text.Source(assets.Source),
				opt.Text.Color(color.White),
				opt.Text.Size(18),
			},
		},
		// Items
		Block: uiex.ClassStyleItem[*ui.Block]{
			"block-image": {
				opt.NoShape(),
				opt.NoEvent(),
			},
		},
		Grid: uiex.ClassStyleItem[*ui.Grid]{
			"grid-header": {
				opt.NoShape(),
			},
			"grid-title": {
				opt.Size(CellWidth*6, CellHeight*3),
				opt.AlignLeft(),
				opt.AlignTop(),
				opt.NoEvent(),
			},
			"grid-decoration": {
				opt.Size(CellWidth*6, CellHeight*3),
				opt.AlignRight(),
				opt.AlignTop(),
				opt.NoEvent(),
			},
			"grid-pad": {
				opt.NoShape(),
				opt.Padding(4),
			},
		},
		ButtonText: uiex.ClassStyleItem[*uiex.ButtonText]{
			"button-pad": {
				opt.Rounding(12.5),
			},
			"button-digit": {
				opt.EventStyle(ui.EventOptions{
					ui.Default:    opt.RGB(59, 59, 59),
					ui.Hover:      opt.RGB(50, 50, 50),
					ui.PressHover: opt.RGB(40, 40, 40),
				}),
			},
			"button-ope": {
				opt.EventStyle(ui.EventOptions{
					ui.Default: opt.RGB(50, 50, 50),
					ui.Hover:   opt.RGB(59, 59, 59),
				}),
			},
			"button-equal": {
				opt.EventStyle(ui.EventOptions{
					ui.Default:    opt.RGB(199, 197, 250),
					ui.Hover:      opt.RGB(182, 180, 228),
					ui.PressHover: opt.RGB(165, 163, 206),
				}),
			},
			"button-minmax": {
				opt.EventStyle(ui.EventOptions{
					ui.Default:    opt.RGB(32, 32, 32),
					ui.Hover:      opt.RGB(50, 50, 50),
					ui.PressHover: opt.RGB(59, 59, 59),
				}),
			},
			"button-close": {
				opt.EventStyle(ui.EventOptions{
					ui.Default:    opt.RGB(32, 32, 32),
					ui.Hover:      opt.RGB(232, 17, 35),
					ui.PressHover: opt.RGB(241, 112, 122),
				}),
			},
		},
		Label: uiex.ClassStyleItem[*uiex.Label]{
			"label-result": {
				opt.PaddingRight(16),
			},
			"label-input": {
				opt.PaddingRight(12),
			},
			"label-metrics": {
				opt.PaddingLeft(12),
			},
			"label-title": {
				opt.NoEvent(),
			},
		},
		// Content
		BasicText: uiex.ClassStyleContent[*uiex.BasicText]{
			"text-title": {
				opt.Text.Size(14),
				opt.Text.AlignLeft(),
				opt.Text.AlignBottom(),
			},
			"text-decoration": {
				opt.Text.Size(26),
			},
			"text-result": {
				opt.Text.AlignRight(),
			},
			"text-input": {
				opt.Text.Size(43),
				opt.Text.AlignRight(),
			},
			"text-metrics": {
				opt.Text.Size(14),
				opt.Text.RGB(32, 196, 0),
				opt.Text.AlignLeft(),
			},
		},
	}
)

func init() {
	// Ebiten icon
	clr := color.RGBA{219, 87, 31, 255}
	vector.DrawFilledRect(ebitenImg, 0, 18-6, 6, 6, clr, false)
	vector.DrawFilledRect(ebitenImg, 3, 18-9, 6, 6, clr, false)
	vector.DrawFilledRect(ebitenImg, 6, 18-12, 6, 6, clr, false)
	vector.DrawFilledRect(ebitenImg, 18-6, 0, 6, 6, clr, false)
	ebitenImg.SubImage(image.Rect(18-3, 0, 18, 3)).(*ebiten.Image).Clear()
}

type padEventHandler struct {
	keys []ebiten.Key
}

func (eh *padEventHandler) State(item ui.Item, is ui.InputState) ui.Event {
	for _, k := range eh.keys {
		switch {
		case is.KeyPressDuration(k) > 0:
			return ui.PressHover
		case is.KeyJustReleased(k):
			return ui.ReleaseHover
		}
	}
	// Use the default mouse event handler after the key bindings checks
	return ui.GetEventHandler().State(item, is)
}

var (
	keyBindings = map[rune][]ebiten.Key{
		'²': {ebiten.KeyApostrophe},
		'÷': {ebiten.KeySlash, ebiten.KeyNumpadDivide},
		'<': {ebiten.KeyBackspace},
		'7': {ebiten.Key7, ebiten.KeyNumpad7},
		'8': {ebiten.Key8, ebiten.KeyNumpad8},
		'9': {ebiten.Key9, ebiten.KeyNumpad9},
		'×': {ebiten.KeyKPMultiply, ebiten.KeyNumpadMultiply},
		'4': {ebiten.Key4, ebiten.KeyNumpad4},
		'5': {ebiten.Key5, ebiten.KeyNumpad5},
		'6': {ebiten.Key6, ebiten.KeyNumpad6},
		'–': {ebiten.KeyMinus, ebiten.KeyNumpadSubtract},
		'1': {ebiten.Key1, ebiten.KeyNumpad1},
		'2': {ebiten.Key2, ebiten.KeyNumpad2},
		'3': {ebiten.Key3, ebiten.KeyNumpad3},
		'+': {ebiten.KeyNumpadAdd},
		'0': {ebiten.Key0, ebiten.KeyNumpad0},
		'.': {ebiten.KeyComma, ebiten.KeyNumpadDecimal},
		'=': {ebiten.KeyEqual, ebiten.KeyNumpadEqual, ebiten.KeyEnter, ebiten.KeyNumpadEnter},
	}
)

func NewButtonPad(content string, calc *app.Calculator, t rune, class string) *uiex.ButtonText {
	return uiex.NewButtonText(content).WithOptions(
		opt.ButtonText.Options(
			opt.Classes("button-pad", class),
			opt.EventHandler(&padEventHandler{
				keys: keyBindings[t],
			}),
			opt.EventAction(ui.EventOptions{
				ui.ReleaseHover: func(_ ui.Item) {
					calc.ProcessToken(t)
				},
			}),
		),
	)
}

type Game struct {
	calc         *app.Calculator
	resultLabel  *uiex.Label
	inputLabel   *uiex.Label
	metricsLabel *uiex.Label
	layout       *ui.Layout
}

var (
	dragged      bool
	dragX, dragY int
)

func New() *Game {
	calc := app.New()

	layout := ui.NewLayout(Columns, Rows, image.Rect(0, 0, CellWidth, CellHeight))
	layout.Grid().SetItemOptions(
		opt.RGB(32, 32, 32),
		// Hack: Remove aliasing on the edges of the undecorated window
		opt.Margin(-1),
	)

	// Window icon + title
	titleBar := ui.NewGrid(3, 2)
	titleBar.SetClasses("grid-header", "grid-title")
	// Icon
	titleBar.Add(0, 0, 1, 1, ui.NewBlock().WithOptions(
		opt.Block.Options(
			opt.Classes("block-image"),
		),
		opt.Block.Content(uiex.NewImage(ebitenImg)),
	))
	// Title
	titleBar.Add(1, 0, 2, 1, uiex.NewLabel("calc.exe").WithOptions(
		opt.Label.Options(
			opt.Classes("label-title"),
		),
		opt.Label.Text(opt.Text.Classes("text-title")),
	))
	// Window top bar decoration
	decorationBar := ui.NewGrid(3, 2)
	decorationBar.SetClasses("grid-header", "grid-decoration")
	// Minimize button
	decorationBar.Add(0, 0, 1, 1, uiex.NewButtonText("–").WithOptions(
		opt.ButtonText.Options(
			opt.Classes("button-minmax"),
			opt.EventAction(ui.EventOptions{
				ui.ReleaseHover: func(_ ui.Item) {
					ebiten.MinimizeWindow()
				},
			}),
		),
		opt.ButtonText.Text(opt.Text.Classes("text-decoration")),
	))
	// Maximize button
	decorationBar.Add(1, 0, 1, 1, uiex.NewButtonText("¤").WithOptions(
		opt.ButtonText.Options(
			opt.Classes("button-minmax"),
			opt.EventAction(ui.EventOptions{
				ui.ReleaseHover: func(_ ui.Item) {
					if ebiten.IsWindowMaximized() {
						ebiten.RestoreWindow()
					} else {
						ebiten.MaximizeWindow()
					}
				},
			}),
		),
		opt.ButtonText.Text(opt.Text.Classes("text-decoration")),
	))
	// Close button
	decorationBar.Add(2, 0, 1, 1, uiex.NewButtonText("×").WithOptions(
		opt.ButtonText.Options(
			opt.Classes("button-close"),
			opt.EventAction(ui.EventOptions{
				ui.ReleaseHover: func(_ ui.Item) {
					os.Exit(0)
				},
			}),
		),
		opt.ButtonText.Text(opt.Text.Classes("text-decoration")),
	))
	// Window header bar
	headerBar := ui.NewGrid(2, 1).WithOptions(
		opt.Grid.Options(
			opt.Classes("grid-header"),
			ui.WithCustomUpdateFunc(func(g *ui.Grid, is ui.InputState) {
				switch {
				case g.Hovered():
					switch is.MouseButtonPressDuration(ebiten.MouseButtonLeft) {
					case 1:
						dragged = !titleBar.Hovered() && !decorationBar.Hovered()
						if dragged {
							dragX, dragY = ebiten.CursorPosition()
						}
					case 0:
						dragged = false
					}
				case g.JustUnhovered():
					// Keep dragging even if not hovering for a moment
					dragged = dragged && is.MouseButtonPressDuration(ebiten.MouseButtonLeft) >= 1
				}
			}),
		),
	)

	// Add title and decoration to header bar
	headerBar.Add(0, 0, 1, 1, titleBar)
	headerBar.Add(1, 0, 1, 1, decorationBar)

	// Top calculation result
	resultLabel := uiex.NewLabel("").WithOptions(
		opt.Label.Options(opt.Classes("label-result")),
		opt.Label.Text(opt.Text.Classes("text-result")),
	)
	// Bottom user input
	inputLabel := uiex.NewLabel("").WithOptions(
		opt.Label.Options(opt.Classes("label-input")),
		opt.Label.Text(opt.Text.Classes("text-input")),
	)
	// Bottom user input
	metricsLabel := uiex.NewLabel("").WithOptions(
		opt.Label.Options(opt.Classes("label-metrics")),
		opt.Label.Text(opt.Text.Classes("text-metrics")),
	)

	// Top operators row
	padGrid := ui.NewGrid(4, 5).WithOptions(
		opt.Grid.Options(
			opt.NoShape(),
			opt.Padding(4),
			opt.PaddingTop(11),
		),
	)
	// Row 1
	padGrid.Add(0, 0, 1, 1, NewButtonPad("n²", calc, '²', "button-ope"))
	padGrid.Add(1, 0, 1, 1, NewButtonPad("n³", calc, '³', "button-ope"))
	padGrid.Add(2, 0, 1, 1, NewButtonPad("÷", calc, '÷', "button-ope"))
	padGrid.Add(3, 0, 1, 1, NewButtonPad("<×", calc, '<', "button-ope"))
	// Row 2
	padGrid.Add(0, 1, 1, 1, NewButtonPad("7", calc, '7', "button-digit"))
	padGrid.Add(1, 1, 1, 1, NewButtonPad("8", calc, '8', "button-digit"))
	padGrid.Add(2, 1, 1, 1, NewButtonPad("9", calc, '9', "button-digit"))
	padGrid.Add(3, 1, 1, 1, NewButtonPad("×", calc, '×', "button-ope"))
	// Row 3
	padGrid.Add(0, 2, 1, 1, NewButtonPad("4", calc, '4', "button-digit"))
	padGrid.Add(1, 2, 1, 1, NewButtonPad("5", calc, '5', "button-digit"))
	padGrid.Add(2, 2, 1, 1, NewButtonPad("6", calc, '6', "button-digit"))
	padGrid.Add(3, 2, 1, 1, NewButtonPad("–", calc, '–', "button-ope"))
	// Row 4
	padGrid.Add(0, 3, 1, 1, NewButtonPad("1", calc, '1', "button-digit"))
	padGrid.Add(1, 3, 1, 1, NewButtonPad("2", calc, '2', "button-digit"))
	padGrid.Add(2, 3, 1, 1, NewButtonPad("3", calc, '3', "button-digit"))
	padGrid.Add(3, 3, 1, 1, NewButtonPad("+", calc, '+', "button-ope"))
	// Row 5
	padGrid.Add(0, 4, 1, 1, NewButtonPad("±", calc, '±', "button-digit"))
	padGrid.Add(1, 4, 1, 1, NewButtonPad("0", calc, '0', "button-digit"))
	padGrid.Add(2, 4, 1, 1, NewButtonPad(",", calc, '.', "button-digit"))
	padGrid.Add(3, 4, 1, 1, NewButtonPad("=", calc, '=', "button-equal"))

	// Add components to the UI layout
	layout.Grid().Add(0, 0, Columns, 3, headerBar)
	layout.Grid().Add(0, 4, Columns, 1, resultLabel)
	layout.Grid().Add(0, 5, Columns, 4, inputLabel)
	layout.Grid().Add(0, 11, Columns, 1, metricsLabel)
	layout.Grid().Add(0, 12, Columns, 13, padGrid)

	// Apply custom theme
	appTheme.Apply(layout)

	return &Game{
		calc:         calc,
		resultLabel:  resultLabel,
		inputLabel:   inputLabel,
		metricsLabel: metricsLabel,
		layout:       layout,
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	// Window dragging
	cx, cy := ebiten.CursorPosition()
	if dragged {
		dx, dy := cx-dragX, cy-dragY
		winX, winY := ebiten.WindowPosition()
		ebiten.SetWindowPosition(winX+dx, winY+dy)
	}

	g.resultLabel.Text().SetText(g.calc.TopString())
	g.inputLabel.Text().SetText(g.calc.InputString())
	inputLen := len(g.inputLabel.Text().Text())
	if inputLen <= 14 {
		g.inputLabel.Text().(*uiex.BasicText).SetSize(43)
	} else {
		g.inputLabel.Text().(*uiex.BasicText).SetSize(43 - 3*float64(inputLen-14))
	}

	g.layout.Update(image.Point{}, ui.GetInputState())

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.metricsLabel.Text().SetText(fmt.Sprintf(
		"FPS: %.2f - TPS: %.2f",
		ebiten.ActualFPS(),
		ebiten.ActualTPS(),
	))

	g.layout.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	g.layout.SetDimensions(outsideWidth, outsideHeight)
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetTPS(TPS)
	ebiten.SetVsyncEnabled(false)
	ebiten.SetFullscreen(false)
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSizeLimits(320, 500, -1, -1)

	g := New()
	if err := ebiten.RunGameWithOptions(g, &ebiten.RunGameOptions{
		GraphicsLibrary: ebiten.GraphicsLibraryOpenGL,
	}); err != nil && err != ebiten.Termination {
		log.Fatal("main: ", err)
	}
}
