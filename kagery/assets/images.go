package assets

import (
	"bytes"
	_ "embed"
	"image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed images/gopher.png
	gopherSrc   []byte
	GopherImage *ebiten.Image

	//go:embed images/gopherbg.png
	gopherBgSrc   []byte
	GopherBgImage *ebiten.Image

	//go:embed images/noise.png
	noiseSrc   []byte
	NoiseImage *ebiten.Image

	//go:embed images/normal.png
	normalSrc   []byte
	NormalImage *ebiten.Image
)

func init() {
	var err error

	img, err := png.Decode(bytes.NewReader(gopherSrc))
	if err != nil {
		log.Fatal(err)
	}
	GopherImage = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(gopherBgSrc))
	if err != nil {
		log.Fatal(err)
	}
	GopherBgImage = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(noiseSrc))
	if err != nil {
		log.Fatal(err)
	}
	NoiseImage = ebiten.NewImageFromImage(img)

	img, err = png.Decode(bytes.NewReader(normalSrc))
	if err != nil {
		log.Fatal(err)
	}
	NormalImage = ebiten.NewImageFromImage(img)
}
