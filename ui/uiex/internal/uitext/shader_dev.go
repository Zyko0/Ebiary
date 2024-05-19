//go:build dev

package uitext

import (
	"log"

	"github.com/Zyko0/Ebiary/asset"
	"github.com/hajimehoshi/ebiten/v2"

	_ "embed"
)

var (
	liveTextShader *asset.LiveAsset[*ebiten.Shader]
)

func init() {
	var err error

	liveTextShader, err = asset.NewLiveAsset[*ebiten.Shader]("./ui/uiex/internal/uitext/shader/text.kage")
	if err != nil {
		log.Fatal("uiex: couldn't compile text shader: ", err)
	}
}

func Shader() *ebiten.Shader {
	return liveTextShader.Value()
}
