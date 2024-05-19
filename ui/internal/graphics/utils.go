package graphics

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	BrushImage  = ebiten.NewImage(1, 1)
	rectIndices = [6]uint16{0, 1, 2, 1, 2, 3}
)

func init() {
	BrushImage = ebiten.NewImage(1, 1)
	BrushImage.Fill(color.White)
}

type RectOpts struct {
	DstX, DstY          float32
	SrcX, SrcY          float32
	DstWidth, DstHeight float32
	SrcWidth, SrcHeight float32
	R, G, B, A          float32
}

func AppendRectVerticesIndices(vertices []ebiten.Vertex, indices []uint16, index int, opts *RectOpts) ([]ebiten.Vertex, []uint16) {
	sx, sy, dx, dy := opts.SrcX, opts.SrcY, opts.DstX, opts.DstY
	sw, sh, dw, dh := opts.SrcWidth, opts.SrcHeight, opts.DstWidth, opts.DstHeight
	r, g, b, a := opts.R, opts.G, opts.B, opts.A
	vertices = append(vertices, []ebiten.Vertex{
		{
			DstX:   dx,
			DstY:   dy,
			SrcX:   sx,
			SrcY:   sy,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   dx + dw,
			DstY:   dy,
			SrcX:   sx + sw,
			SrcY:   sy,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   dx,
			DstY:   dy + dh,
			SrcX:   sx,
			SrcY:   sy + sh,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		{
			DstX:   dx + dw,
			DstY:   dy + dh,
			SrcX:   sx + sw,
			SrcY:   sy + sh,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
	}...)

	indiceCursor := uint16(index * 4)
	indices = append(indices, []uint16{
		rectIndices[0] + indiceCursor,
		rectIndices[1] + indiceCursor,
		rectIndices[2] + indiceCursor,
		rectIndices[3] + indiceCursor,
		rectIndices[4] + indiceCursor,
		rectIndices[5] + indiceCursor,
	}...)

	return vertices, indices
}

func ColorAsFloat32RGB(clr color.Color) float32 {
	if clr == nil {
		return 0
	}
	r, g, b, _ := clr.RGBA()
	return float32((r&255)<<16 + (g&255)<<8 + b&255)
}
