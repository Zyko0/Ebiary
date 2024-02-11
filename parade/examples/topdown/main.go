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
	WallsHeight = 128
	WallWidth   = 68
)

type Game struct {
	renderer *parade.Renderer

	player *game.Player

	boxEntitiesDepth *ebiten.Image
	boxEntitiesColor *ebiten.Image
	entitiesDepth    *ebiten.Image
	entitiesColor    *ebiten.Image
}

func NewGame() *Game {
	return &Game{
		renderer: parade.NewRenderer(ScreenWidth, ScreenHeight, 1000, parade.Backward),

		player: game.NewPlayer(),

		boxEntitiesDepth: ebiten.NewImage(LayerWidth, LayerHeight),
		boxEntitiesColor: ebiten.NewImage(LayerWidth, LayerHeight),
		entitiesDepth:    ebiten.NewImage(LayerWidth, LayerHeight),
		entitiesColor:    ebiten.NewImage(LayerWidth, LayerHeight),
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
	g.player.Update()
	if g.player.X+game.PlayerSize/2 > LayerWidth-WallWidth {
		g.player.X = LayerWidth - WallWidth - game.PlayerSize/2
	}
	if g.player.X-game.PlayerSize/2 < WallWidth {
		g.player.X = WallWidth + game.PlayerSize/2
	}
	if g.player.Y+game.PlayerSize/2 > LayerHeight-WallWidth {
		g.player.Y = LayerHeight - WallWidth - game.PlayerSize/2
	}
	if g.player.Y-game.PlayerSize/2 < WallWidth {
		g.player.Y = WallWidth + game.PlayerSize/2
	}
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
	// Dynamic entities layer
	c := game.Crate{
		WallWidth + game.CrateSize/2, WallWidth + game.CrateSize/2,
	}
	// Crates
	g.boxEntitiesDepth.Clear()
	c.DrawDepth(g.boxEntitiesDepth)
	g.boxEntitiesColor.Clear()
	c.DrawColor(g.boxEntitiesColor)
	// Player
	g.entitiesDepth.Clear()
	g.player.DrawDepth(g.entitiesDepth)
	g.entitiesColor.Clear()
	g.player.DrawColor(g.entitiesColor)
	// Layers
	g.renderer.DrawLayers(screen.SubImage(image.Rect(0, 0, 1024, 1024)).(*ebiten.Image),
		[]*parade.Layer{
			{
				Z:       WallsHeight,
				Depth:   10,
				Diffuse: game.ImageFloorColor,
				Height:  game.ImageFloorDepth,
			},
			{
				Z:         WallsHeight,
				Depth:     WallsHeight,
				Diffuse:   game.ImageRocksColor,
				Height:    game.ImageWallsDepth,
				BoxMapped: true,
			},
			{
				Z:         WallsHeight,
				Depth:     game.CrateSize,
				Diffuse:   g.boxEntitiesColor,
				Height:    g.boxEntitiesDepth,
				BoxMapped: true,
			},
			{
				Z:       WallsHeight,
				Depth:   game.PlayerHeight,
				Diffuse: g.entitiesColor,
				Height:  g.entitiesDepth,
			},
			{
				Z:         25,
				Depth:     25,
				Diffuse:   game.ImageRocksColor,
				Height:    game.ImageRoofDepth,
				BoxMapped: true,
			},
			/*
				{
					Z:       WallsHeight - 10,
					Depth:   10,
					Diffuse: game.ImageFloorColor,
					Height:  game.ImageFloorDepth,
				},
				{
					Z:         0,
					Depth:     WallsHeight,
					Diffuse:   game.ImageRocksColor,
					Height:    game.ImageWallsDepth,
					BoxMapped: true,
				},
				{
					Z:         WallsHeight - game.CrateSize,
					Depth:     WallsHeight,
					Diffuse:   g.boxEntitiesColor,
					Height:    g.boxEntitiesDepth,
					BoxMapped: true,
				},
				{
					Z:       WallsHeight - game.PlayerHeight,
					Depth:   game.PlayerHeight,
					Diffuse: g.entitiesColor,
					Height:  g.entitiesDepth,
				},
				{
					Z:         0,
					Depth:     25,
					Diffuse:   game.ImageRocksColor,
					Height:    game.ImageRoofDepth,
					BoxMapped: true,
				},*/
			/*{
				Z:       0,
				Depth:   10,
				Diffuse: game.ImageFloorColor,
				Height:  game.ImageFloorDepth,
			},
			{
				Z:         0,
				Depth:     WallsHeight,
				Diffuse:   game.ImageRocksColor,
				Height:    game.ImageWallsDepth,
				BoxMapped: true,
			},
			{
				Z:         0,
				Depth:     25,
				Diffuse:   game.ImageRocksColor,
				Height:    game.ImageRoofDepth,
				BoxMapped: true,
			},
			{
				Z:         WallsHeight + game.CrateSize,
				Depth:     game.CrateSize,
				Diffuse:   g.boxEntitiesColor,
				Height:    g.boxEntitiesDepth,
				BoxMapped: true,
			},
			{
				Z:       WallsHeight + game.PlayerHeight,
				Depth:   game.PlayerHeight,
				Diffuse: g.entitiesColor,
				Height:  g.entitiesDepth,
			},*/
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
		fmt.Sprintf("FPS: %.2f - X: %.02f - Y: %.02f",
			ebiten.ActualFPS(), g.player.X, g.player.Y,
		),
	)
}

func (g *Game) Layout(_, _ int) (int, int) {
	return 1024, 1024 //ScreenWidth, ScreenHeight
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
