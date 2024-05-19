//go:build !dev

package uitext

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	_ "embed"
)

var (
	//go:embed shader/text.kage
	srcTextShader []byte
	textShader    *ebiten.Shader
)

func init() {
	var err error

	textShader, err = ebiten.NewShader(srcTextShader)
	if err != nil {
		log.Fatal("uiex: couldn't compile text shader: ", err)
	}
}

func Shader() *ebiten.Shader {
	return textShader
}
