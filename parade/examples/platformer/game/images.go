package game

import (
	"bytes"
	_ "embed"

	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	// Background

	//go:embed images/ebiten_pillars_depth.png
	srcLayerPillarsDepth   []byte
	ImageLayerPillarsDepth *ebiten.Image

	//go:embed images/ebiten_pillars_color.png
	srcLayerPillarsColor   []byte
	ImageLayerPillarsColor *ebiten.Image

	// Sprites

	//go:embed images/playeridle0.png
	srcPlayerIdle0   []byte
	ImagePlayerIdle0 *ebiten.Image

	//go:embed images/playeridle1.png
	srcPlayerIdle1   []byte
	ImagePlayerIdle1 *ebiten.Image

	//go:embed images/playerjump.png
	srcPlayerJump   []byte
	ImagePlayerJump *ebiten.Image

	//go:embed images/playerrun0.png
	srcPlayerRun0   []byte
	ImagePlayerRun0 *ebiten.Image

	//go:embed images/playerrun1.png
	srcPlayerRun1   []byte
	ImagePlayerRun1 *ebiten.Image
)

func init() {
	// Background

	img, err := png.Decode(bytes.NewReader(srcLayerPillarsDepth))
	if err != nil {
		log.Fatal(err)
	}
	ImageLayerPillarsDepth = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(srcLayerPillarsColor))
	if err != nil {
		log.Fatal(err)
	}
	ImageLayerPillarsColor = ebiten.NewImageFromImage(img)

	// Sprites

	img, err = png.Decode(bytes.NewReader(srcPlayerIdle0))
	if err != nil {
		log.Fatal(err)
	}
	ImagePlayerIdle0 = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(srcPlayerIdle1))
	if err != nil {
		log.Fatal(err)
	}
	ImagePlayerIdle1 = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(srcPlayerJump))
	if err != nil {
		log.Fatal(err)
	}
	ImagePlayerJump = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(srcPlayerRun0))
	if err != nil {
		log.Fatal(err)
	}
	ImagePlayerRun0 = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(srcPlayerRun1))
	if err != nil {
		log.Fatal(err)
	}
	ImagePlayerRun1 = ebiten.NewImageFromImage(img)
}
