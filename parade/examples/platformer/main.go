package main

import (
	"fmt"
	"log"

	"github.com/Zyko0/Ebiary/parade"
	"github.com/Zyko0/Ebiary/parade/examples/platformer/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	TPS          = 60
	ScreenWidth  = 1280
	ScreenHeight = 720

	LayerWidth  = 2560
	LayerHeight = 1440
)

type Game struct {
	renderer *parade.Renderer
	player   *game.Player
}

func NewGame() *Game {
	return &Game{
		renderer: parade.NewRenderer(ScreenWidth, ScreenHeight, 1000, parade.Forward),
		player:   game.NewPlayer(),
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
	py := g.player.Y
	g.player.Update()
	// Floor collision
	if g.player.Y <= 0 && (g.player.State == game.StateFalling || g.player.State == game.StateJumping) {
		g.player.Y = 0
		g.player.Grounded = true
		g.player.VelocityY = 0
		g.player.State, g.player.StateTick = game.StateIdle, 0
	}
	// Find nearest platform
	nearest := game.NearestPlatform(g.player)
	if !nearest.Nil() {
		column := g.player.X > float64(nearest.X-game.PlatformWidth/2)
		column = column && g.player.X < float64(nearest.X+game.PlatformWidth/2)
		landed := !g.player.Grounded && g.player.State == game.StateFalling
		landed = landed && g.player.Y <= float64(nearest.Y-game.PlatformHeight/2)
		landed = landed && py >= float64(nearest.Y-game.PlatformHeight/2)
		landed = landed && column
		if landed {
			g.player.Grounded = true
			g.player.Y = float64(nearest.Y - game.PlatformHeight/2)
			g.player.VelocityY = 0
			g.player.State, g.player.StateTick = game.StateIdle, 0
		}
		if g.player.Y > 0 && !column {
			if g.player.State != game.StateFalling && g.player.State != game.StateJumping {
				g.player.Grounded = false
				g.player.State, g.player.StateTick = game.StateFalling, 0
			}
		}
	}
	/*if !g.nearestPlatform.Nil() {
		landed := !g.player.Grounded && g.player.State == game.StateFalling
		landed = landed && g.player.Y < float64(g.nearestPlatform.Y-game.PlatformHeight/2)
		landed = landed && g.player.X > float64(g.nearestPlatform.X-game.PlatformWidth/2)
		landed = landed && g.player.X < float64(g.nearestPlatform.X+game.PlatformWidth/2)
		if landed {
			g.player.Y = float64(g.nearestPlatform.Y - game.PlatformHeight/2)
			g.player.VelocityY = 0
			g.player.State, g.player.StateTick = game.StateIdle, 0
		}
	}*/
	//fmt.Println("nearest", nearest)
	// Camera matrices update
	camY := LayerHeight/2 - g.player.Y
	g.renderer.Camera().SetPosition(g.player.X, camY, 0)
	g.renderer.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	//screen.DrawImage(assets.ImageLayer, nil)

	g.renderer.DrawLayers(screen,
		[]*parade.Layer{
			{
				Z:       500,
				Depth:   250,
				Diffuse: game.ImageLayerPillarsColor,
				Height:  game.ImageLayerPillarsDepth,
			},
			{
				Z:       50,
				Depth:   100,
				Diffuse: game.PlatformLayerDepthImage,
				Height:  game.PlatformLayerDepthImage,
			},
		},
		&parade.DrawLayersOptions{
			/*OffsetX:      0,
			OffsetY:      0,*/
			//Antialiasing: true,
		},
	)
	// Draw player
	g.player.Draw(screen)
	//screen.DrawImage(game.PlatformLayerDepthImage, nil)

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("FPS: %.2f - grounded: %v - x %.02f y %.02f",
			ebiten.ActualFPS(), g.player.Grounded, g.player.X, g.player.Y,
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
		GraphicsLibrary: ebiten.GraphicsLibraryOpenGL,
	}); err != ebiten.Termination {
		log.Fatalf("error at RunGame: %v", err)
	}
}
