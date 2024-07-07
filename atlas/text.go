package atlas

import (
	"image"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type textAtlas struct {
	*Atlas
	glyphs map[rune]*Image
}

var (
	textAtlases  = map[text.Face]*textAtlas{}
	textDrawList = &DrawList{}
)

func TextAtlasImage() *ebiten.Image {
	for _, v := range textAtlases {
		return v.native
	}
	return nil
}

func DrawText(dst *ebiten.Image, str string, face text.Face, opts *text.DrawOptions) {
	atlas, ok := textAtlases[face]
	if !ok {
		atlas = &textAtlas{
			Atlas: New(1024, 1024, &NewAtlasOptions{
				// TODO: Figure with metrics()
				// TODO: Not sure we can assume anything here, but should be okay
				// as it's only for first writes
				MinSize: image.Pt(1, 1),
			}),
			glyphs: map[rune]*Image{},
		}
		textAtlases[face] = atlas
	}

	glyphs := text.AppendGlyphs(nil, str, face, &opts.LayoutOptions)
	str = strings.ReplaceAll(str, "\n", "")
	dc := &DrawCommand{}
	for i, r := range str {
		glyph := glyphs[i]
		if glyph.Image == nil {
			continue
		}
		/*
			// text/v2.Draw version
			geom := opts.GeoM
			geom.Translate(glyph.X, glyph.Y)
			dst.DrawImage(glyph.Image, &ebiten.DrawImageOptions{
				GeoM: geom,
			})
			continue*/

		cached, ok := atlas.glyphs[r]
		if !ok {
			cached = atlas.NewImage(
				glyph.Image.Bounds().Dx(),
				glyph.Image.Bounds().Dy(),
			)
			// TODO: below can be stacked as triangles also
			// and executed before the final flush of drawlist
			geom := ebiten.GeoM{}
			geom.Translate(
				float64(cached.bounds.Min.X),
				float64(cached.bounds.Min.Y),
			)
			atlas.native.DrawImage(glyph.Image, &ebiten.DrawImageOptions{
				GeoM: geom,
			})

			atlas.glyphs[r] = cached
		}
		// Add glyph draw command
		dc.Image = cached
		dc.GeoM.Reset()
		dc.GeoM.Concat(opts.GeoM)
		dc.GeoM.Translate(glyph.X, glyph.Y)
		textDrawList.Add(dc)
	}
	// Flush to destination
	textDrawList.Flush(dst, &DrawOptions{})

	//fmt.Println("lencached", len(atlas.glyphs))
}
