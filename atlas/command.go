package atlas

import (
	"github.com/Zyko0/Ebiary/atlas/internal/graphics"
	"github.com/hajimehoshi/ebiten/v2"
)

type drawRange struct {
	atlas *Atlas
	end   int
}

type DrawList struct {
	ranges []drawRange
	vx     []ebiten.Vertex
	ix     []uint16
}

type DrawCommand struct {
	Image *Image

	ColorScale ebiten.ColorScale
	GeoM       ebiten.GeoM
}

func (dl *DrawList) Add(commands ...*DrawCommand) {
	if len(commands) == 0 {
		return
	}

	var batch *drawRange

	if len(dl.ranges) > 0 {
		batch = &dl.ranges[len(dl.ranges)-1]
	} else {
		dl.ranges = append(dl.ranges, drawRange{
			atlas: commands[0].Image.atlas,
		})
		batch = &dl.ranges[0]
	}
	// Add vertices and indices
	opts := &graphics.RectOpts{}
	for _, cmd := range commands {
		if cmd.Image.atlas != batch.atlas {
			dl.ranges = append(dl.ranges, drawRange{
				atlas: cmd.Image.atlas,
			})
			batch = &dl.ranges[len(dl.ranges)-1]
		}

		x, y := cmd.GeoM.Apply(0, 0)
		w, h := cmd.GeoM.Apply(
			float64(cmd.Image.bounds.Dx()),
			float64(cmd.Image.bounds.Dy()),
		)
		w -= x
		h -= y
		/*fmt.Printf("xy: %.0f; %.0f - wh: %.0f; %.0f (wh %d, %d)\n",
			x, y, w, h, cmd.Image.bounds.Dx(), cmd.Image.bounds.Dy(),
		)*/
		opts.R = cmd.ColorScale.R()
		opts.G = cmd.ColorScale.G()
		opts.B = cmd.ColorScale.B()
		opts.A = cmd.ColorScale.A()
		opts.SrcX = float32(cmd.Image.bounds.Min.X)
		opts.SrcY = float32(cmd.Image.bounds.Min.Y)
		opts.SrcWidth = float32(cmd.Image.bounds.Dx())
		opts.SrcHeight = float32(cmd.Image.bounds.Dy())
		opts.DstX = float32(x)
		opts.DstY = float32(y)
		opts.DstWidth = float32(w)
		opts.DstHeight = float32(h)
		dl.vx, dl.ix = graphics.AppendRectVerticesIndices(
			dl.vx, dl.ix, batch.end, opts,
		)

		batch.end++
	}
}

type DrawOptions struct {
	ColorScaleMode ebiten.ColorScaleMode
	Blend          ebiten.Blend
	Filter         ebiten.Filter
	Address        ebiten.Address
	AntiAlias      bool
}

func (dl *DrawList) Flush(dst *ebiten.Image, opts *DrawOptions) {
	index := 0
	for _, r := range dl.ranges {
		dst.DrawTriangles(
			dl.vx[index*4:r.end*4],
			dl.ix[index*6:r.end*6],
			r.atlas.native,
			&ebiten.DrawTrianglesOptions{
				ColorScaleMode: opts.ColorScaleMode,
				Blend:          opts.Blend,
				Filter:         opts.Filter,
				Address:        opts.Address,
				AntiAlias:      opts.AntiAlias,
			},
		)
		index += r.end
	}
	dl.ranges = dl.ranges[:0]
	dl.vx = dl.vx[:0]
	dl.ix = dl.ix[:0]
}
