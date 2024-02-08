package main

import (
	"fmt"
	"image"
	"log"

	"github.com/Zyko0/Ebiary/parade"
	"github.com/Zyko0/Ebiary/parade/examples/topdown/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	TPS          = 60
	ScreenWidth  = 1280
	ScreenHeight = 1280

	LayerWidth  = 1024
	LayerHeight = 1024
)

type Game struct {
	renderer *parade.Renderer

	entitiesImage *ebiten.Image
}

func NewGame() *Game {
	return &Game{
		renderer: parade.NewRenderer(ScreenWidth, ScreenHeight, 50),
	}
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}
	// Player update
	//TODO:
	// Floor collision
	//TODO:
	// Find nearest platform

	// Camera matrices update
	/*camY := LayerHeight/2 - g.player.Y
	g.camera.SetPosition(g.player.X, camY, 0)*/
	g.renderer.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//screen.DrawImage(game.ImageRocksColor, nil)

	g.renderer.DrawLayers(screen.SubImage(image.Rect(0, 0, 1024, 1024)).(*ebiten.Image),
		[]*parade.Layer{
			{
				Z:       0,
				Depth:   5,
				Diffuse: game.ImageFloorColor,
				Height:  game.ImageFloorDepth,
			},
			{
				Z:         0,
				Depth:     50,
				Diffuse:   game.ImageRocksColor,
				Height:    game.ImageWallsDepth,
				BoxMapped: true,
			},
			{
				Z:         0,
				Depth:     10,
				Diffuse:   game.ImageRocksColor,
				Height:    game.ImageRoofDepth,
				BoxMapped: true,
			},
		},
		&parade.DrawLayersOptions{
			/*OffsetX:      0,
			OffsetY:      0,*/
			//Antialiasing: true,
		},
	)
	// Draw player
	//g.player.Draw(screen) // TODO:
	//screen.DrawImage(game.PlatformLayerDepthImage, nil)

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("FPS: %.2f",
			ebiten.ActualFPS(),
		),
	)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func main() {
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetFullscreen(false)
	ebiten.SetVsyncEnabled(false) // TODO: true

	if err := ebiten.RunGameWithOptions(NewGame(), &ebiten.RunGameOptions{
		//GraphicsLibrary: ebiten.GraphicsLibraryOpenGL,
	}); err != ebiten.Termination {
		log.Fatalf("error at RunGame: %v", err)
	}
}
