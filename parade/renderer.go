package parade

import (
	"github.com/Zyko0/Ebiary/parade/internal/camera"
	"github.com/Zyko0/Ebiary/parade/internal/utils"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/hajimehoshi/ebiten/v2"
)

type Renderer struct {
	width  float64
	height float64
	depth  float64
	camera *camera.CameraOrthographic
}

func NewRenderer(width, height, depth int) *Renderer {
	w, h, d := float64(width), float64(height), float64(depth)

	return &Renderer{
		width:  w,
		height: h,
		depth:  d,
		camera: camera.NewOrthographic(w, h, d),
	}
}

func (r *Renderer) Update() {
	r.camera.Update()
}

func (r *Renderer) Layout() (int, int) {
	return int(r.width), int(r.height)
}

func (r *Renderer) Camera() Camera {
	return r.camera
}

type DrawLayersOptions struct {
	// Antialiasing uses the native AntiAlias draw option from ebitengine.
	Antialiasing bool
}

func (r *Renderer) DrawLayers(screen *ebiten.Image, layers []*Layer, opts *DrawLayersOptions) {
	proj := r.camera.ProjectionMatrix()
	view := r.camera.ViewMatrix()
	pvinv := mgl64.Mat4(proj).Mul4(mgl64.Mat4(view)).Inv()
	x, y, z := r.camera.Position()
	screenSize := screen.Bounds()
	/*sort.SliceStable(layers, func(i int, j int) bool {
		return layers[i].Z > layers[j].Z
	})*/
	println("x, y", screen.Bounds().Dx(), screen.Bounds().Dy())
	for _, l := range layers {
		var boxmapped float32

		if l.BoxMapped {
			boxmapped = 1
		}
		lb := l.Height.Bounds()
		vertices, indices := utils.AppendRectVerticesIndices(
			nil, nil, 0, &utils.RectOpts{
				DstX:      0,
				DstY:      0,
				SrcX:      0,
				SrcY:      0,
				DstWidth:  float32(screen.Bounds().Dx()),
				DstHeight: float32(screen.Bounds().Dy()),
				SrcWidth:  float32(r.width),  //l.Height.Bounds().Dx()),
				SrcHeight: float32(r.height), //l.Height.Bounds().Dy()),
				R:         float32(l.Z),
				G:         float32(l.Depth),
				B:         0,
				A:         0,
			},
		)
		screen.DrawTrianglesShader(vertices, indices, liveShader.Value(), &ebiten.DrawTrianglesShaderOptions{
			Uniforms: map[string]any{
				"Normal":   float32(0),
				"Specular": float32(0),

				"ScreenSize": []float32{
					float32(screenSize.Dx()), float32(screenSize.Dy()),
				},
				"LayerSize": []float32{
					float32(lb.Dx()), float32(lb.Dy()),
				},
				"BoxMapping": boxmapped,

				"CameraPVMatrixInv": pvinv[:],
				"CameraPosition": []float64{
					x, y, z,
				},
				"LightPosition": []float32{
					0, 0, -2,
				},
			},
			Images: [4]*ebiten.Image{
				l.Height,
				l.Diffuse,
			},
			AntiAlias: opts.Antialiasing,
		})
	}
}
