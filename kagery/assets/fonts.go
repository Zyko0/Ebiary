package assets

import (
	"bytes"
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	//go:embed fonts/noto-sans-regular.ttf
	notoTTF    []byte
	NotoSource *text.GoTextFaceSource

	//go:embed fonts/sourcecodepro-regular.ttf
	sourcecodeproTTF    []byte
	SourceCodeProSource *text.GoTextFaceSource
)

func init() {
	var err error

	NotoSource, err = text.NewGoTextFaceSource(bytes.NewReader(notoTTF))
	if err != nil {
		log.Fatal("couldn't create face source from ttf: ", err)
	}

	SourceCodeProSource, err = text.NewGoTextFaceSource(bytes.NewReader(sourcecodeproTTF))
	if err != nil {
		log.Fatal("couldn't create face source from ttf: ", err)
	}
}
