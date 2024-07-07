package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/Zyko0/Ebiary/atlas"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	ScreenWidth  = 1920
	ScreenHeight = 1080
)

var (
	source *text.GoTextFaceSource
	font   *text.GoTextFace
)

func init() {
	var err error

	source, err = text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatal("couldn't load font: ", err)
	}
	font = &text.GoTextFace{
		Source: source,
		Size:   12,
	}
}

type mode bool

const (
	DrawImage mode = false
	DrawList  mode = true
)

func (m mode) String() string {
	if m == DrawImage {
		return "DrawImage"
	}
	return "DrawList"
}

type App struct {
	mode mode
	str  string
}

func New() *App {
	var str string
	for i := 0; i < 20; i++ {
		str += "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.\n"
		str += "Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.\n"
		str += "Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.\n"
		str += "Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.\n"
	}

	return &App{
		str: str,
	}
}

func (a *App) Update() error {
	// Mode switch
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		a.mode = !a.mode
	}

	return nil
}

func (a *App) DrawImages(screen *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Translate(0, 32)
	text.Draw(screen, a.str, font, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: geom,
		},
		LayoutOptions: text.LayoutOptions{
			LineSpacing: font.Size,
		},
	})
}

func (a *App) DrawList(screen *ebiten.Image) {
	geom := ebiten.GeoM{}
	geom.Translate(0, 32)
	atlas.DrawText(screen, a.str, font, &text.DrawOptions{
		DrawImageOptions: ebiten.DrawImageOptions{
			GeoM: geom,
		},
		LayoutOptions: text.LayoutOptions{
			LineSpacing: font.Size,
		},
	})
}

func (a *App) Draw(screen *ebiten.Image) {
	switch a.mode {
	case DrawImage:
		a.DrawImages(screen)
	case DrawList:
		a.DrawList(screen)
	}

	if img := atlas.TextAtlasImage(); img != nil {
		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(800, 32)
		screen.DrawImage(img, opts)
	}

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("Mode: %v - Characters: %d - FPS: %.02f",
			a.mode,
			len(a.str),
			ebiten.ActualFPS(),
		),
	)
}

func (a *App) Layout(ow, oh int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetCursorShape(ebiten.CursorShapeCrosshair)
	ebiten.SetVsyncEnabled(false)

	/*f, err := os.Create("beat.prof")
	if err != nil {
		log.Fatal(err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal(err)
	}
	defer pprof.StopCPUProfile()*/

	if err := ebiten.RunGame(New()); err != nil {
		log.Fatal("err run game: ", err)
	}
}
