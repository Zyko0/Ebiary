package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	PlatformWidth  = 320
	PlatformHeight = 64

	PlatformXSpacing   = PlatformWidth * 2
	PlatformYSpacing   = PlatformHeight * 6
	PlatformOddOffsetY = PlatformYSpacing / 2
)

type Platform struct {
	X, Y int
}

func (p Platform) Over(player *Player) bool {
	return player.Y >= float64(p.Y+PlatformHeight/2)
}

func (p Platform) Nil() bool {
	return p.X == 0 && p.Y == 0
}

func NearestPlatform(p *Player) Platform {
	x, y := p.X, p.Y
	x = x / PlatformXSpacing
	y = y / PlatformYSpacing
	signx := 1.
	if x < 0 {
		signx = -1
	}
	yoff := 0.
	if int(math.Abs(x+0.5*signx))%2 == 1 {
		yoff = 0.5
	}
	return Platform{
		X: int(x+0.5*signx)*PlatformXSpacing + PlatformWidth/2,
		Y: int(y+0.5)*PlatformYSpacing + int(yoff*PlatformYSpacing), // + PlatformHeight/2,
	}
}

var (
	whiteImg *ebiten.Image

	PlatformLayerDepthImage *ebiten.Image
	PlatformLayerColorImage *ebiten.Image
)

func init() {
	whiteImg = ebiten.NewImage(1, 1)
	whiteImg.Fill(color.White)
	// Generate platforms
	PlatformLayerDepthImage = ebiten.NewImage(2560, 1440)
	PlatformLayerColorImage = ebiten.NewImage(2560, 1440)
	opts := &ebiten.DrawImageOptions{}
	for i := 0; i < 2; i++ {
		for x := 0.; x < 2560/PlatformXSpacing; x++ {
			for y := 0.; y < 1440/PlatformYSpacing; y++ {
				opts.GeoM.Reset()
				opts.GeoM.Scale(PlatformWidth, PlatformHeight)
				yoff := 0.
				if int(x)%2 == 0 {
					yoff = -PlatformYSpacing / 2
				}
				opts.GeoM.Translate(
					x*PlatformXSpacing,
					1440-y*PlatformYSpacing+yoff+PlatformHeight,
				)
				PlatformLayerDepthImage.DrawImage(whiteImg, opts)
			}
		}
	}
}
