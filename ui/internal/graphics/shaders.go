//go:build !dev

package graphics

import (
	_ "embed"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	//go:embed shaders/item.kage
	srcItemShader []byte
	itemShader    *ebiten.Shader
)

func init() {
	var err error

	itemShader, err = ebiten.NewShader(srcItemShader)
	if err != nil {
		log.Fatal("ui: couldn't compile item shader: ", err)
	}
}

func ItemShader() *ebiten.Shader {
	return itemShader
}
