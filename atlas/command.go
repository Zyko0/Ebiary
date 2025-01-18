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
	const (
		// 16384 commands by batch (65536 indices max)
		maxBatchEnd = 65536 / 4
	)

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
		if cmd.Image.atlas != batch.atlas || batch.end >= maxBatchEnd {
			dl.ranges = append(dl.ranges, drawRange{
				atlas: cmd.Image.atlas,
			})
			batch = &dl.ranges[len(dl.ranges)-1]
		}

		// Dst attributes
		bounds := cmd.Image.bounds
		opts.Dsts[0] = graphics.Pt(cmd.GeoM.Apply(0, 0))
		opts.Dsts[1] = graphics.Pt(cmd.GeoM.Apply(
			float64(bounds.Dx()), 0,
		))
		opts.Dsts[2] = graphics.Pt(cmd.GeoM.Apply(
			0, float64(bounds.Dy()),
		))
		opts.Dsts[3] = graphics.Pt(cmd.GeoM.Apply(
			float64(bounds.Dx()), float64(bounds.Dy()),
		))

		// Color and source attributes
		opts.R = cmd.ColorScale.R()
		opts.G = cmd.ColorScale.G()
		opts.B = cmd.ColorScale.B()
		opts.A = cmd.ColorScale.A()
		opts.SrcX0 = float32(bounds.Min.X)
		opts.SrcY0 = float32(bounds.Min.Y)
		opts.SrcX1 = float32(bounds.Max.X)
		opts.SrcY1 = float32(bounds.Max.Y)

		dl.vx, dl.ix = graphics.AppendRectVerticesIndices(
			dl.vx, dl.ix, batch.end, opts,
		)

		//if len(dl.ix) >=

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
			dl.vx[index*4:(index+r.end)*4],
			dl.ix[index*6:(index+r.end)*6],
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

type DrawShaderOptions struct {
	Blend       ebiten.Blend
	Uniforms    map[string]any
	ExtraImages [3]*ebiten.Image
	AntiAlias   bool
}

func (dl *DrawList) FlushWithShader(dst *ebiten.Image, shader *ebiten.Shader, opts *DrawShaderOptions) {
	topts := &ebiten.DrawTrianglesShaderOptions{}
	if opts != nil {
		topts.Blend = opts.Blend
		topts.Uniforms = opts.Uniforms
		topts.Images = [4]*ebiten.Image{
			nil,
			opts.ExtraImages[0],
			opts.ExtraImages[1],
			opts.ExtraImages[2],
		}
		topts.AntiAlias = opts.AntiAlias
	}
	index := 0
	for _, r := range dl.ranges {
		topts.Images[0] = r.atlas.native
		dst.DrawTrianglesShader(
			dl.vx[index*4:(index+r.end)*4],
			dl.ix[index*6:(index+r.end)*6],
			shader,
			topts,
		)
		index += r.end
	}
	// Clear buffers
	dl.ranges = dl.ranges[:0]
	dl.vx = dl.vx[:0]
	dl.ix = dl.ix[:0]
}
