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

	//go:embed images/rocks_color.png
	srcRocksColor   []byte
	ImageRocksColor *ebiten.Image

	//go:embed images/floor_depth.png
	srcFloorDepth   []byte
	ImageFloorDepth *ebiten.Image

	//go:embed images/floor_color.png
	srcFloorColor   []byte
	ImageFloorColor *ebiten.Image

	//go:embed images/walls_depth.png
	srcWallsDepth   []byte
	ImageWallsDepth *ebiten.Image

	//go:embed images/roof_depth.png
	srcRoofDepth   []byte
	ImageRoofDepth *ebiten.Image
)

func init() {
	// Rocks color

	img, err := png.Decode(bytes.NewReader(srcRocksColor))
	if err != nil {
		log.Fatal(err)
	}
	ImageRocksColor = ebiten.NewImageFromImage(img)

	// Floor

	img, err = png.Decode(bytes.NewReader(srcFloorDepth))
	if err != nil {
		log.Fatal(err)
	}
	ImageFloorDepth = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(srcFloorColor))
	if err != nil {
		log.Fatal(err)
	}
	ImageFloorColor = ebiten.NewImageFromImage(img)

	// Walls

	img, err = png.Decode(bytes.NewReader(srcWallsDepth))
	if err != nil {
		log.Fatal(err)
	}
	ImageWallsDepth = ebiten.NewImageFromImage(img)

	// Roof

	img, err = png.Decode(bytes.NewReader(srcRoofDepth))
	if err != nil {
		log.Fatal(err)
	}
	ImageRoofDepth = ebiten.NewImageFromImage(img)

	// Sprites

}
