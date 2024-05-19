package graphics

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Pipeline struct {
	bounds image.Rectangle
	layers []*layer
}

func NewPipeline(bounds image.Rectangle) *Pipeline {
	return &Pipeline{
		bounds: bounds,
		layers: []*layer{
			{},
		},
	}
}

func (pp *Pipeline) Clear() {
	for _, l := range pp.layers {
		l.Clear()
	}
}

func (pp *Pipeline) EnsureLayers(z int) {
	if z > len(pp.layers) {
		panic("graphics: missing layer for item's z index")
	}
	if z == len(pp.layers) {
		pp.layers = append(pp.layers, &layer{})
	}
}

func (pp *Pipeline) Add(item *ItemPrimitive, area image.Rectangle) {
	pp.EnsureLayers(item.Z)
	pp.layers[item.Z].Add(item, area)
}

func (pp *Pipeline) AddContent(c *ContentPrimitive) {
	pp.EnsureLayers(c.Z)
	pp.layers[c.Z].AddContent(c)
}

// TODO: Add image
// TODO: Add draw func

func (pp *Pipeline) Draw(screen *ebiten.Image, offset image.Point) {
	// Draw layers sorted by Z
	for _, l := range pp.layers {
		l.Draw(screen, offset)
	}
}
