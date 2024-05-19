package graphics

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	buffersCount = 64
)

type layer struct {
	// Primitives
	index        int
	uniformData  [][]float32
	uniformExtra [][]float32
	vx           [][]ebiten.Vertex
	ix           [][]uint16
	// Custom content
	contents []*ContentPrimitive
}

func (l *layer) Clear() {
	// Note: no need to clear uniforms
	l.index = 0
	for i := range l.vx {
		l.vx[i] = l.vx[i][:0]
	}
	for i := range l.ix {
		l.ix[i] = l.ix[i][:0]
	}
	l.contents = l.contents[:0]
}

func shapeArgs(s Shape, x, y float32) (float32, float32) {
	switch s {
	case ShapeBox:
		return x, y
	case ShapeEllipse:
		return x, y
	case ShapeRhombus:
		return x, y
	case ShapeOctogon:
		return x, y
	default:
		return 0, 0
	}
}

func (l *layer) Add(item *ItemPrimitive, area image.Rectangle) {
	buffer := l.index / buffersCount
	index := l.index % buffersCount
	// Grow slices if necessary
	if len(l.vx) == buffer {
		l.uniformData = append(l.uniformData, make([]float32, buffersCount*4*2))
		l.uniformExtra = append(l.uniformExtra, make([]float32, buffersCount*2))
		l.vx = append(l.vx, []ebiten.Vertex{})
		l.ix = append(l.ix, []uint16{})
	}
	// X, Y and size normalization
	width := float32(area.Dx()) - item.MarginLeft - item.MarginRight
	height := float32(area.Dy()) - item.MarginTop - item.MarginBottom
	x := width / max(width, height)
	y := height / max(width, height)
	// Append uniforms data
	l.uniformData[buffer][index*8+0] = float32(item.Shape)
	l.uniformData[buffer][index*8+1] = float32(item.Rounding) / max(width, height)
	l.uniformData[buffer][index*8+2] = float32(item.BorderWidth) / max(width, height)
	l.uniformData[buffer][index*8+3] = float32(item.ColorMinFactor)
	l.uniformData[buffer][index*8+4] = float32(item.ColorFilling)
	arg0, arg1 := shapeArgs(item.Shape, x, y)
	l.uniformData[buffer][index*8+5] = arg0
	l.uniformData[buffer][index*8+6] = arg1
	l.uniformData[buffer][index*8+7] = item.ColorAlpha
	// Append uniforms extra data
	l.uniformExtra[buffer][index*2+0] = AAFactor / max(width, height)
	l.uniformExtra[buffer][index*2+1] = 0
	// Append geometry
	l.vx[buffer], l.ix[buffer] = AppendRectVerticesIndices(
		l.vx[buffer], l.ix[buffer], index, &RectOpts{
			DstX:      float32(area.Min.X) + item.MarginLeft,
			DstY:      float32(area.Min.Y) + item.MarginTop,
			DstWidth:  width,
			DstHeight: height,
			SrcX:      -x,
			SrcY:      -y,
			SrcWidth:  x * 2,
			SrcHeight: y * 2,
			R:         float32(index),
			G:         ColorAsFloat32RGB(item.ColorMin),
			B:         ColorAsFloat32RGB(item.ColorMax),
			A:         ColorAsFloat32RGB(item.BorderColor),
		},
	)
	l.index++
}

func (l *layer) AddContent(c *ContentPrimitive) {
	l.contents = append(l.contents, c)
}

func (l *layer) Draw(dst *ebiten.Image, offset image.Point) {
	/*if err := ItemShader.Error(); err != nil {//TODO: re-enable this in dev mode
		fmt.Println("err:", err)
		return
	}*/
	// Draw primitives
	for i := range l.ix {
		for j := range l.vx[i] {
			l.vx[i][j].DstX += float32(offset.X)
			l.vx[i][j].DstY += float32(offset.Y)
		}
		dst.DrawTrianglesShader(l.vx[i], l.ix[i], ItemShader(), &ebiten.DrawTrianglesShaderOptions{
			Uniforms: map[string]any{
				"GammaExp": GammaExp,
				"Data":     l.uniformData[i],
				"Extra":    l.uniformExtra[i],
			},
			Blend: ebiten.BlendSourceOver,
		})
	}
	// Draw custom content
	aa := int(math.Round(AAFactor + 0.499999))
	paa := image.Point{aa, aa}
	for _, c := range l.contents {
		clip := c.Clip
		clip.Min, clip.Max = c.Clip.Min.Add(paa), c.Clip.Max.Sub(paa)
		clip.Min, clip.Max = clip.Min.Add(offset), clip.Max.Add(offset)
		c.Content.Draw(dst.SubImage(clip).(*ebiten.Image))
	}
}
