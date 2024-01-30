package jfa

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type JFA struct {
	rect       image.Rectangle
	img0, img1 *ebiten.Image
}

func New(rect image.Rectangle) *JFA {
	return &JFA{
		rect: rect,
		img0: ebiten.NewImageWithOptions(rect, nil),
		img1: ebiten.NewImageWithOptions(rect, nil),
	}
}

type Encoding byte

const (
	EncodingDistanceGreyscale Encoding = iota
	EncodingUV
	/*EncodingRGB768
	EncodingRGBA1024*/
)

type DistanceType byte

const (
	DistanceExterior DistanceType = iota
	DistanceInterior
)

type ColorMask byte

const (
	ColorMaskAlpha ColorMask = iota
	ColorMaskGreyscale
	ColorMaskR
	ColorMaskG
	ColorMaskB
)

type GenerateOptions struct {
	// PlainValueOptions defines a list of constraints of minimum
	// value thresholds (in 0-1 range) to consider a pixel's color as a "plain" value
	// (as opposed to empty).
	// For example ColorMaskGreyScale with value: 0.1 means that
	// the pixel's greyscale value must be over 0.1 to contribute to a plain value.
	// Multiple options in the list act as an "AND" operation.
	// By default, it is considered that a pixel is a "plain" value if its
	// alpha channel > 0.
	PlainValueThresholds map[ColorMask]float64
	// DistanceType defines whether the resulting distance encoding should be a:
	// - Exterior distance to compute the exterior minimal distance to the shape.
	// - Interior distance to compute the interior minimal distance to the shape's edges
	// By default, the exterior distance is encoded.
	DistanceType DistanceType
	// EdgesPlain defines whether or not to consider the image's boundaries (or
	// edges of the image) as plain value for distance calculation.
	// By default, it is false and edges do not contribute to plain values.
	EdgesPlain bool
	// Encoding specifies the way to encode the final distance image.
	// By default, the distance will be encoded as greyscale with alpha = 255.
	Encoding Encoding
	// EncodingScale specifies a value to mutliply the resulting color by.
	// By default, the scale is 1 (it can exceed 1).
	EncodingScale float64
	// Steps defines the number of iterations to perform, a higher value means
	// an higher output quality.
	// Default is 255 steps.
	Steps int
	// JumpDistance defines the initial jump distance in the JFA.
	// Default is 8.
	JumpDistance int
}

func (jfa *JFA) Generate(dst, src *ebiten.Image, opts *GenerateOptions) {
	jfab := jfa.img0.Bounds()
	srcb := src.Bounds()
	dstb := dst.Bounds()
	if srcb.Dx() != jfab.Dx() || srcb.Dy() != jfab.Dy() {
		panic("jfa: source image should have the same bounds")
	}
	if dstb.Dx() != jfab.Dx() || dstb.Dy() != jfab.Dy() {
		panic("jfa: destination image should have the same bounds")
	}
	// Handle default values
	if opts == nil {
		opts = &GenerateOptions{
			EncodingScale: 1,
			Steps:         255,
			JumpDistance:  8,
		}
	}
	if len(opts.PlainValueThresholds) == 0 {
		opts.PlainValueThresholds[ColorMaskAlpha] = 0.01
	}
	if opts.EncodingScale == 0 {
		opts.EncodingScale = 1
	}
	if opts.Steps == 0 {
		opts.Steps = 255
	}
	if opts.JumpDistance == 0 {
		opts.JumpDistance = 8
	}

	shader := JFAShader(opts)
	blend := ebiten.BlendCopy
	vertices := []ebiten.Vertex{
		{
			DstX: 0,
			DstY: 0,
			SrcX: float32(srcb.Min.X),
			SrcY: float32(srcb.Min.Y),
		},
		{
			DstX: float32(jfab.Dx()),
			DstY: 0,
			SrcX: float32(srcb.Max.X),
			SrcY: float32(srcb.Min.Y),
		},
		{
			DstX: 0,
			DstY: float32(jfab.Dy()),
			SrcX: float32(srcb.Min.X),
			SrcY: float32(srcb.Max.Y),
		},
		{
			DstX: float32(jfab.Dx()),
			DstY: float32(jfab.Dy()),
			SrcX: float32(srcb.Max.X),
			SrcY: float32(srcb.Max.Y),
		},
	}
	indices := []uint16{0, 1, 2, 1, 2, 3}
	jfa.img0.Clear()
	jfa.img1.Clear()
	// Initial pass
	jfa.img0.DrawTrianglesShader(vertices, indices, shader, &ebiten.DrawTrianglesShaderOptions{
		Uniforms: map[string]any{
			"ColorMaskAlpha":     opts.PlainValueThresholds[ColorMaskAlpha],
			"ColorMaskGreyscale": opts.PlainValueThresholds[ColorMaskGreyscale],
			"ColorMaskR":         opts.PlainValueThresholds[ColorMaskR],
			"ColorMaskG":         opts.PlainValueThresholds[ColorMaskG],
			"ColorMaskB":         opts.PlainValueThresholds[ColorMaskB],
			"FirstPass":          float32(1),
		},
		Images: [4]*ebiten.Image{
			src,
		},
		Blend: blend,
	})
	// Next passes
	s, d := jfa.img0, jfa.img1
	jd := opts.JumpDistance
	for i := 0; i < opts.Steps; i++ {
		d.DrawTrianglesShader(vertices, indices, shader, &ebiten.DrawTrianglesShaderOptions{
			Uniforms: map[string]any{
				"FirstPass":    float32(0),
				"JumpDistance": float32(jd),
			},
			Images: [4]*ebiten.Image{
				s,
			},
			Blend: blend,
		})
		jd = max(jd/2, 1)
		s, d = d, s
	}
	// Encoding
	for i := range vertices {
		vertices[i].DstX += float32(dstb.Min.X)
		vertices[i].DstY += float32(dstb.Min.Y)
	}
	dst.DrawTrianglesShader(vertices, indices, EncodingShader(), &ebiten.DrawTrianglesShaderOptions{
		Uniforms: map[string]any{
			"Encoding": float32(opts.Encoding),
			"Scale":    float32(opts.EncodingScale),
		},
		Images: [4]*ebiten.Image{
			s,
		},
		Blend: blend,
	})
}

func (jfa *JFA) Bounds() image.Rectangle {
	return jfa.rect
}
