package assets

import (
	"bytes"
	_ "embed"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	srcPlayer   []byte
	ImagePlayer *ebiten.Image

	//go:embed images/layer3.png
	srcLayer   []byte
	ImageLayer *ebiten.Image

	//go:embed images/layer_pillars_back.png
	srcLayerPillarsBack   []byte
	ImageLayerPillarsBack *ebiten.Image

	//go:embed images/layer_pillars_front.png
	srcLayerPillarsFront   []byte
	ImageLayerPillarsFront *ebiten.Image
)

func init() {
	img, err := png.Decode(bytes.NewReader(srcLayerPillarsBack))
	if err != nil {
		log.Fatal(err)
	}
	ImageLayerPillarsBack = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(srcLayerPillarsFront))
	if err != nil {
		log.Fatal(err)
	}
	ImageLayerPillarsFront = ebiten.NewImageFromImage(img)

	// FIXME: this one is debug
	img, err = png.Decode(bytes.NewReader(srcLayer))
	if err != nil {
		log.Fatal(err)
	}
	ImageLayer = ebiten.NewImageFromImage(img)
}
