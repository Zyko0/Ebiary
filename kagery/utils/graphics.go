package utils

import "github.com/hajimehoshi/ebiten/v2"

var (
	rectIndices = [6]uint16{0, 1, 2, 1, 2, 3}
)

type RectOpts struct {
	DstX, DstY          float32
	SrcX, SrcY          float32
	DstWidth, DstHeight float32
	SrcWidth, SrcHeight float32
	R, G, B, A          float32
}

func AppendRectVerticesIndices(vertices []ebiten.Vertex, indices []uint16, opts *RectOpts) ([]ebiten.Vertex, []uint16) {
	sx, sy, dx, dy := opts.SrcX, opts.SrcY, opts.DstX, opts.DstY
	sw, sh, dw, dh := opts.SrcWidth, opts.SrcHeight, opts.DstWidth, opts.DstHeight
	r, g, b, a := opts.R, opts.G, opts.B, opts.A
	vertices = append(vertices,
		ebiten.Vertex{
			DstX:   dx,
			DstY:   dy,
			SrcX:   sx,
			SrcY:   sy,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		ebiten.Vertex{
			DstX:   dx + dw,
			DstY:   dy,
			SrcX:   sx + sw,
			SrcY:   sy,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		ebiten.Vertex{
			DstX:   dx,
			DstY:   dy + dh,
			SrcX:   sx,
			SrcY:   sy + sh,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		ebiten.Vertex{
			DstX:   dx + dw,
			DstY:   dy + dh,
			SrcX:   sx + sw,
			SrcY:   sy + sh,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
	)

	indices = append(indices,
		rectIndices[0],
		rectIndices[1],
		rectIndices[2],
		rectIndices[3],
		rectIndices[4],
		rectIndices[5],
	)

	return vertices, indices
}
