package game

import (
	"bytes"
	_ "embed"
	"image/color"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	ImageWhite *ebiten.Image

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

	//go:embed images/crate_color.png
	srcCrateColor   []byte
	ImageCrateColor *ebiten.Image

	//go:embed images/player_depth.png
	srcPlayerDepth   []byte
	ImagePlayerDepth *ebiten.Image

	//go:embed images/player_color.png
	srcPlayerColor   []byte
	ImagePlayerColor *ebiten.Image
)

func init() {
	ImageWhite = ebiten.NewImage(1, 1)
	ImageWhite.Fill(color.White)

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

	// Entities

	img, err = png.Decode(bytes.NewReader(srcCrateColor))
	if err != nil {
		log.Fatal(err)
	}
	ImageCrateColor = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(srcPlayerDepth))
	if err != nil {
		log.Fatal(err)
	}
	ImagePlayerDepth = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(srcPlayerColor))
	if err != nil {
		log.Fatal(err)
	}
	ImagePlayerColor = ebiten.NewImageFromImage(img)
}
