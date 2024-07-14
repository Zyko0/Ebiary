package atlas

import (
	"github.com/Zyko0/Ebiary/atlas/internal/graphics"
	"github.com/hajimehoshi/ebiten/v2"
)

type drawRange struct {
	atlas *Atlas
	end   int
}

// DrawList stores triangle versions of DrawImage calls when
// all images are sub-images of an atlas.
// Temporary vertices and indices can be re-used after calling
// Flush, so it is more efficient to keep a reference to a DrawList
// instead of creating a new one every frame.
type DrawList struct {
	ranges []drawRange
	vx     []ebiten.Vertex
	ix     []uint16
}

// DrawCommand is the equivalent of the regular DrawImageOptions
// but only including options that will not break batching.
// Filter, Address, Blend and AntiAlias are determined at Flush()
type DrawCommand struct {
	Image *Image

	ColorScale ebiten.ColorScale
	GeoM       ebiten.GeoM
}

// Add adds DrawImage commands to the DrawList, images from multiple
// atlases can be added but they will break the previous batch bound to
// a different atlas, requiring an additional draw call internally.
// So, it is better to have the maximum of consecutive DrawCommand images
// sharing the same atlas.
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
		opts.DstWidth = float32(w - x)
		opts.DstHeight = float32(h - y)
		dl.vx, dl.ix = graphics.AppendRectVerticesIndices(
			dl.vx, dl.ix, batch.end, opts,
		)

		batch.end++
	}
}

// DrawOptions are additional options that will be applied to
// all draw commands from the draw list when calling Flush().
type DrawOptions struct {
	ColorScaleMode ebiten.ColorScaleMode
	Blend          ebiten.Blend
	Filter         ebiten.Filter
	Address        ebiten.Address
	AntiAlias      bool
}

// Flush executes all the draw commands as the smallest possible
// amount of draw calls, and then clears the list for next uses.
func (dl *DrawList) Flush(dst *ebiten.Image, opts *DrawOptions) {
	var topts *ebiten.DrawTrianglesOptions
	if opts != nil {
		topts = &ebiten.DrawTrianglesOptions{
			ColorScaleMode: opts.ColorScaleMode,
			Blend:          opts.Blend,
			Filter:         opts.Filter,
			Address:        opts.Address,
			AntiAlias:      opts.AntiAlias,
		}
	}
	index := 0
	for _, r := range dl.ranges {
		dst.DrawTriangles(
			dl.vx[index*4:r.end*4],
			dl.ix[index*6:r.end*6],
			r.atlas.native,
			topts,
		)
		index += r.end
	}
	// Clear buffers
	dl.ranges = dl.ranges[:0]
	dl.vx = dl.vx[:0]
	dl.ix = dl.ix[:0]
}
