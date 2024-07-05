package main

import (
	"image"
	"image/color"
	"log"
	"math/rand"

	"github.com/Zyko0/Ebiary/atlas"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 1920
	ScreenHeight = 1080
)

type App struct {
	atlas *atlas.Atlas
}

func New() *App {
	return &App{
		atlas: atlas.New(1024, 1024, &atlas.NewAtlasOptions{
			MinSize: image.Pt(8, 8),
		}),
	}
}

func (a *App) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		w, h := 8+rand.Intn(16), 8+rand.Intn(16)
		img := a.atlas.NewImage(w, h)
		if img != nil {
			img.Image().Fill(color.RGBA{
				R: uint8(rand.Intn(255)),
				G: uint8(rand.Intn(255)),
				B: uint8(rand.Intn(255)),
				A: 255,
			})
		}
	}
	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	screen.DrawImage(a.atlas.Image(), &ebiten.DrawImageOptions{
		Blend: ebiten.BlendCopy,
	})
}

func (a *App) Layout(ow, oh int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)

	if err := ebiten.RunGame(New()); err != nil {
		log.Fatal("err run game: ", err)
	}
}
