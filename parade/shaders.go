package parade

import (
	_ "embed"
	"log"

	"github.com/Zyko0/Ebiary/asset"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	shaders map[string]*ebiten.Shader
	//ESSgo:embed shaders/projection.tmpl.kage
	//srcLayerShaderTmpl []byte
	//ShaderLayer        *ebiten.Shader
	liveShader *asset.LiveAsset[*ebiten.Shader]
)

func ShaderLayer() *ebiten.Shader {
	return liveShader.Value()
}

func init() {
	var err error

	liveShader, err = asset.NewLiveAsset[*ebiten.Shader]("shaders/projection.tmpl.kage")
	if err != nil {
		log.Fatal("shaders:", err)
	}

	shaders = map[string]*ebiten.Shader{}
}