package graphics

import "github.com/hajimehoshi/ebiten/v2"

var rectIndices = [6]uint16{0, 1, 2, 1, 2, 3}

type RectOpts struct {
	DstX, DstY          float32
	SrcX, SrcY          float32
	DstWidth, DstHeight float32
	SrcWidth, SrcHeight float32
	R, G, B, A          float32
}

// adjustDestinationPixel is the original ebitengine implementation found here:
// https://github.com/hajimehoshi/ebiten/blob/v2.8.0-alpha.1/internal/graphics/vertex.go#L102-L126
func adjustDestinationPixel(x float32) float32 {
	// Avoid the center of the pixel, which is problematic (#929, #1171).
	// Instead, align the vertices with about 1/3 pixels.
	//
	// The intention here is roughly this code:
	//
	//     float32(math.Floor((float64(x)+1.0/6.0)*3) / 3)
	//
	// The actual implementation is more optimized than the above implementation.
	ix := float32(int(x))
	if x < 0 && x != ix {
		ix -= 1
	}
	frac := x - ix
	switch {
	case frac < 3.0/16.0:
		return ix
	case frac < 8.0/16.0:
		return ix + 5.0/16.0
	case frac < 13.0/16.0:
		return ix + 11.0/16.0
	default:
		return ix + 16.0/16.0
	}
}

func AppendRectVerticesIndices(vertices []ebiten.Vertex, indices []uint16, index int, opts *RectOpts) ([]ebiten.Vertex, []uint16) {
	sx, sy, dx, dy := opts.SrcX, opts.SrcY, opts.DstX, opts.DstY
	sw, sh, dw, dh := opts.SrcWidth, opts.SrcHeight, opts.DstWidth, opts.DstHeight
	r, g, b, a := opts.R, opts.G, opts.B, opts.A
	vertices = append(vertices,
		ebiten.Vertex{
			DstX:   adjustDestinationPixel(dx),
			DstY:   adjustDestinationPixel(dy),
			SrcX:   sx,
			SrcY:   sy,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		ebiten.Vertex{
			DstX:   adjustDestinationPixel(dx + dw),
			DstY:   adjustDestinationPixel(dy),
			SrcX:   sx + sw,
			SrcY:   sy,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		ebiten.Vertex{
			DstX:   adjustDestinationPixel(dx),
			DstY:   adjustDestinationPixel(dy + dh),
			SrcX:   sx,
			SrcY:   sy + sh,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
		ebiten.Vertex{
			DstX:   adjustDestinationPixel(dx + dw),
			DstY:   adjustDestinationPixel(dy + dh),
			SrcX:   sx + sw,
			SrcY:   sy + sh,
			ColorR: r,
			ColorG: g,
			ColorB: b,
			ColorA: a,
		},
	)

	indiceCursor := uint16(index * 4)
	indices = append(indices,
		rectIndices[0]+indiceCursor,
		rectIndices[1]+indiceCursor,
		rectIndices[2]+indiceCursor,
		rectIndices[3]+indiceCursor,
		rectIndices[4]+indiceCursor,
		rectIndices[5]+indiceCursor,
	)

	return vertices, indices
}
