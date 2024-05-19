package assets

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	//go:embed fonts/noto-sans-regular.ttf
	ttf    []byte
	Source *text.GoTextFaceSource
)

func init() {
	var err error

	Source, err = text.NewGoTextFaceSource(bytes.NewReader(ttf))
	if err != nil {
		log.Fatal("couldn't create face source from ttf: ", err)
	}
}
