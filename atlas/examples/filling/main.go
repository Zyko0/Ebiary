package main

import (
	"image"
	"image/color"
	"log"
	"math/rand"
	"slices"

	"github.com/Zyko0/Ebiary/atlas"
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	ScreenWidth  = 1920
	ScreenHeight = 1080
)

type App struct {
	atlas *atlas.Atlas
	imgs  []*atlas.Image
}

func New() *App {
	return &App{
		atlas: atlas.New(1024, 1024, &atlas.NewAtlasOptions{
			MinSize: image.Pt(16, 16),
		}),
	}
}

func (a *App) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		//if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		for i := 0; i < 8; i++ {
			w, h := 16+rand.Intn(48), 16+rand.Intn(48)
			//w, h = 32, 32
			img := a.atlas.NewImage(w, h)
			if img != nil {
				img.Image().Fill(color.RGBA{
					R: uint8(rand.Intn(255)),
					G: uint8(rand.Intn(255)),
					B: uint8(rand.Intn(255)),
					A: 255,
				})
				a.imgs = append(a.imgs, img)
			}
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		for i := 0; i < 8; i++ {
			if len(a.imgs) > 0 {
				idx := rand.Intn(len(a.imgs))
				img := a.imgs[idx]
				a.imgs = slices.Delete(a.imgs, idx, idx+1)
				a.atlas.Free(img)
			}
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
