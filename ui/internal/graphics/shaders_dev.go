//go:build dev

package graphics

import (
	_ "embed"
	"log"

	"github.com/Zyko0/Ebiary/asset"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	itemShader    *asset.LiveAsset[*ebiten.Shader]
)

func init() {
	var err error

	itemShader, err = asset.NewLiveAsset[*ebiten.Shader]("./ui/internal/graphics/shaders/item.kage")
	if err != nil {
		log.Fatal("ui: couldn't compile item shader: ", err)
	}
}

func ItemShader() *ebiten.Shader {
	return itemShader.Value()
}