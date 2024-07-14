package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/Zyko0/Ebiary/atlas"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	ScreenWidth  = 1920
	ScreenHeight = 1080
)

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

type Command struct {
	index int
	geom  ebiten.GeoM
}

type App struct {
	mode        mode
	atlas       *atlas.Atlas
	atlasImages []*atlas.Image
	drawList    *atlas.DrawList

	subImages []*ebiten.Image

	commands []Command
}

func New() *App {
	app := &App{
		atlas: atlas.New(1024, 1024, &atlas.NewAtlasOptions{
			MinSize: image.Pt(16, 16),
		}),
		drawList: &atlas.DrawList{},
	}
	for {
		w, h := 16+rand.Intn(48), 16+rand.Intn(48)
		img := app.atlas.NewImage(w, h)
		if img == nil {
			break
		}
		if img != nil {
			img.Image().Fill(color.RGBA{
				R: uint8(rand.Intn(255)),
				G: uint8(rand.Intn(255)),
				B: uint8(rand.Intn(255)),
				A: 255,
			})
			app.atlasImages = append(app.atlasImages, img)
			app.subImages = append(app.subImages, img.Image())
		}
	}

	return app
}

func (a *App) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		for i := 0; i < 50; i++ {
			geom := ebiten.GeoM{}
			geom.Scale(0.5+rand.Float64(), 0.5+rand.Float64())
			geom.Translate(
				rand.Float64()*ScreenWidth,
				rand.Float64()*ScreenHeight,
			)
			index := rand.Intn(len(a.atlasImages))
			a.commands = append(a.commands, Command{
				index: index,
				geom:  geom,
			})
		}
	}
	// Mode switch
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		a.mode = !a.mode
	}

	return nil
}

func (a *App) DrawImages(screen *ebiten.Image) {
	if len(a.subImages) > 0 {
		opts := &ebiten.DrawImageOptions{}
		for _, c := range a.commands {
			opts.GeoM = c.geom
			screen.DrawImage(a.subImages[c.index], opts)
		}
	}
}

func (a *App) DrawList(screen *ebiten.Image) {
	if len(a.atlasImages) > 0 {
		dc := &atlas.DrawCommand{}
		for _, c := range a.commands {
			dc.Image = a.atlasImages[c.index]
			dc.GeoM = c.geom
			a.drawList.Add(dc)
		}
		// Flush
		a.drawList.Flush(screen, &atlas.DrawOptions{})
	}
}

func (a *App) Draw(screen *ebiten.Image) {
	switch a.mode {
	case DrawImage:
		a.DrawImages(screen)
	case DrawList:
		a.DrawList(screen)
	}

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("Mode: %v - Commands: %d - FPS: %.02f",
			a.mode,
			len(a.commands),
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

	if err := ebiten.RunGame(New()); err != nil {
		log.Fatal("err run game: ", err)
	}
}
