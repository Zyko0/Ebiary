package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/Zyko0/Ebiary/ui"

	"github.com/Zyko0/Ebiary/ui/opt"
	_ "github.com/Zyko0/Ebiary/ui/theme"
	"github.com/Zyko0/Ebiary/ui/uiex"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	TPS          = 60
	ScreenWidth  = 1920
	ScreenHeight = 1080
	Columns      = 40
	Rows         = 30
	CellSize     = 32
)

type Game struct {
	layout *ui.Layout
}

var (
	img = ebiten.NewImage(4000, 256)
)

func init() {
	img.Fill(color.RGBA{0, 0, 128, 255})
	clr := color.RGBA{255, 255, 255, 255}
	for x := float32(128); x < 4000; x += 256 {
		vector.DrawFilledCircle(
			img, x, 128, 128, clr, false,
		)
		clr.R -= 15
		clr.G -= 15
		clr.B -= 15
	}
	//rand.Seed(42)
}

const str string = `Hello how are you doing? I'm doing fine
Hello how are you doing? I'm doing fine
Hello how are you doing? I'm doing fine
Hello how are you doing? I'm doing fine
Hello how are you doing? I'm doing fine
Hello how are you doing? I'm doing fine
Hello how are you doing? I'm doing fine
Hello how are you doing? I'm doing fine
Hello how are you doing? I'm doing fine
`

func New() *Game {
	layout := ui.NewLayout(Columns, Rows, image.Rect(0, 0, CellSize, CellSize))
	layout.Grid().SetItemOptions(
		opt.RGB(0, 128, 0),
	)

	menu := uiex.NewMenuBar(Columns)
	menu.Add(0, 0, 2, 1, uiex.NewButtonText("File"))
	menu.Add(2, 0, 2, 1, uiex.NewButtonText("Settings"))

	p := uiex.NewTextInput().WithOptions(
		opt.TextInput.Options(
			opt.RGB(192, 192, 192),
			opt.Decorations(
				uiex.NewScrollbar().WithOptions(
					opt.Scrollbar.Direction(uiex.DirectionHorizontal),
					// Some space for the second scrollbar overlap
					opt.Scrollbar.Options(opt.PaddingRight(24)),
				),
				uiex.NewScrollbar().WithOptions(
					opt.Scrollbar.Direction(uiex.DirectionVertical),
					// Some space for the second scrollbar overlap
					opt.Scrollbar.Options(opt.PaddingBottom(24)),
				),
			),
		),
		opt.TextInput.RichText(),
	)

	p.Text().SetText(str)
	//txt := p.Text().(*uiex.RichText)
	//txt.PushColorFg(color.RGBA{0, 255, 0, 255})
	//txt.PushColorBg(color.RGBA{0, 0, 255, 255})
	//txt.PushBold()
	//txt.PushItalic()
	//txt.Append("Hey you how old are you?")
	//txt.Pop()
	//txt.Append(" You don't have to answer, but I'd prefer if you do")
	//txt.PushColorBg(color.RGBA{255, 0, 0, 255})
	/*txt.Append(" are ")
	txt.Append("you!")*/

	layout.Grid().Add(0, 0, Columns, 1, menu)
	layout.Grid().Add(0, 1, Columns, Rows-1, p)

	return &Game{
		layout: layout,
	}
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	g.layout.Update(image.Point{64, 64}, ui.GetInputState())

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{128, 0, 128, 255})
	g.layout.Draw(screen)

	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()))
}

func (g *Game) Layout(_, _ int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetTPS(TPS)
	ebiten.SetVsyncEnabled(false)
	ebiten.SetFullscreen(false)
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)

	g := New()
	if err := ebiten.RunGameWithOptions(g, &ebiten.RunGameOptions{
		GraphicsLibrary: ebiten.GraphicsLibraryOpenGL,
	}); err != nil && err != ebiten.Termination {
		log.Fatal("main: ", err)
	}
}
