package main

import (
	"fmt"
	"log"

	"github.com/Zyko0/Ebiary/parade"
	"github.com/Zyko0/Ebiary/parade/examples/topdown/game"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	TPS          = 60
	ScreenWidth  = 1024 //1280
	ScreenHeight = 1024 //1280

	LayerWidth  = 1024
	LayerHeight = 1024
	WallsHeight = 128
	WallWidth   = 68
)

type Game struct {
	renderer *parade.Renderer

	player *game.Player
	crate  *game.Crate

	roomX, roomY int

	boxEntitiesDepth *ebiten.Image
	boxEntitiesColor *ebiten.Image
	entitiesDepth    *ebiten.Image
	entitiesColor    *ebiten.Image
}

func NewGame() *Game {
	return &Game{
		renderer: parade.NewRenderer(ScreenWidth, ScreenHeight, 1000, parade.Backward),

		player: game.NewPlayer(),
		crate: &game.Crate{
			X: LayerWidth / 2,
			Y: LayerHeight / 2,
		},

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
	x, y := g.player.X, g.player.Y
	g.player.Update()
	// Room border collisions
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
	dx, dy := g.player.X-x, g.player.Y-y
	// Crates collisions
	if g.crate.Pushed(g.player) {
		// Push the crate in the player's direction
		g.crate.X += dx
		g.crate.Y += dy
		// Check if crates collide with a wall
		var cx, cy float64
		if g.crate.X+game.CrateSize/2 > LayerWidth-WallWidth {
			cx = g.crate.X
			g.crate.X = LayerWidth - WallWidth - game.CrateSize/2
			cx -= g.crate.X
		}
		if g.crate.X-game.CrateSize/2 < WallWidth {
			cx = g.crate.X
			g.crate.X = WallWidth + game.CrateSize/2
			cx -= g.crate.X
		}
		if g.crate.Y+game.CrateSize/2 > LayerHeight-WallWidth {
			cy = g.crate.Y
			g.crate.Y = LayerHeight - WallWidth - game.CrateSize/2
			cy -= g.crate.Y
		}
		if g.crate.Y-game.CrateSize/2 < WallWidth {
			cy = g.crate.Y
			g.crate.Y = WallWidth + game.CrateSize/2
			cy -= g.crate.Y
		}
		g.player.X -= cx
		g.player.Y -= cy
	}
	// Doors mechanisms
	for _, d := range game.Doors {
		if d.DistanceTo(g.player) < game.DoorTriggerDistance {
			d.TriggerActivation(true)
		} else {
			d.TriggerActivation(false)
		}
		d.Update()
	}

	g.renderer.Camera().SetPosition(0, 0, -512) //512, 512, -512)
	g.renderer.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Crates
	g.boxEntitiesDepth.Clear()
	for _, d := range game.Doors {
		d.DrawDepth(g.boxEntitiesDepth)
	}
	g.crate.DrawDepth(g.boxEntitiesDepth)
	g.boxEntitiesColor.Clear()
	for _, d := range game.Doors {
		d.DrawColor(g.boxEntitiesColor)
	}
	g.crate.DrawColor(g.boxEntitiesColor)
	// Player
	g.entitiesDepth.Clear()
	g.player.DrawDepth(g.entitiesDepth)
	g.entitiesColor.Clear()
	g.player.DrawColor(g.entitiesColor)
	// Layers
	g.renderer.DrawLayers(screen, []*parade.Layer{
		// Floor
		{
			Z:       WallsHeight,
			Depth:   10,
			Diffuse: game.ImageFloorColor,
			Height:  game.ImageFloorDepth,
		},
		// Room walls
		{
			Z:         WallsHeight,
			Depth:     WallsHeight,
			Diffuse:   game.ImageRocksColor,
			Height:    game.ImageWallsDepth,
			BoxMapped: true,
		},
		// Crates and doors
		{
			Z:         WallsHeight,
			Depth:     WallsHeight,
			Diffuse:   g.boxEntitiesColor,
			Height:    g.boxEntitiesDepth,
			BoxMapped: true,
		},
		// Player
		{
			Z:       WallsHeight,
			Depth:   game.PlayerHeight,
			Diffuse: g.entitiesColor,
			Height:  g.entitiesDepth,
		},
		// Walls roof
		{
			Z:         25,
			Depth:     25,
			Diffuse:   game.ImageRocksColor,
			Height:    game.ImageRoofDepth,
			BoxMapped: true,
		},
	}, nil,
	)

	ebitenutil.DebugPrint(screen,
		fmt.Sprintf("FPS: %.2f - X: %.02f - Y: %.02f",
			ebiten.ActualFPS(), g.player.X, g.player.Y,
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
