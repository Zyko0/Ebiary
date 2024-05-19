package ui

import (
	"image"
	"math"
	"slices"

	"github.com/Zyko0/Ebiary/ui/internal/graphics"
)

type Grid struct {
	*itemImpl

	width  int
	height int

	grid  [][]Item
	items []Item
}

func NewGrid(columns, rows int) *Grid {
	grid := make([][]Item, rows)
	for y := range grid {
		grid[y] = make([]Item, columns)
	}
	g := &Grid{
		width:  columns,
		height: rows,

		grid:  grid,
		items: make([]Item, 0),
	}
	g.itemImpl = newItem(g)

	return g
}

func (g *Grid) update(c *context, area image.Rectangle, z int) {
	if g.Skipped() {
		return
	}
	// Grid
	clamped, inner := g.adjustInnerArea(area)
	g.itemImpl.geom.lastFullRegion = inner
	g.itemImpl.geom.lastRegion = clamped
	g.itemImpl.geom.lastCursor = c.Cursor
	g.itemImpl.update(c, g, clamped, z)
	// Decorations
	for _, d := range g.itemImpl.decorations {
		if d.Visible() {
			d.update(c, clamped, z)
		}
	}
	// Children
	unitX, unitY := clamped.Dx()/g.width, clamped.Dy()/g.height
	/*remX := float64(inner.Dx()%g.width) / float64(g.width)
	remY := float64(inner.Dy()%g.height) / float64(g.height)*/
	itemAreas := map[Item]image.Rectangle{}
	for y, row := range g.grid {
		for x, i := range row {
			// Skip empty cells
			if i == nil {
				continue
			}
			// Process cell's item area
			if _, ok := itemAreas[i]; !ok {
				itemAreas[i] = image.Rectangle{
					Min: image.Pt(math.MaxInt32, math.MaxInt32),
					Max: image.Pt(math.MinInt32, math.MinInt32),
				}
			}
			// Padding offsets
			/*xoff := int(float64(x) * remX)
			yoff := int(float64(y) * remY)*/
			r := itemAreas[i]
			r.Min.X = min(r.Min.X, clamped.Min.X+x*unitX)
			r.Min.Y = min(r.Min.Y, clamped.Min.Y+y*unitY)
			r.Max.X = max(r.Max.X, clamped.Min.X+(x+1)*unitX)
			r.Max.Y = max(r.Max.Y, clamped.Min.Y+(y+1)*unitY)
			itemAreas[i] = r
		}
	}
	// Update children base items and their own children
	for i, r := range itemAreas {
		i.update(c, r, z+1)
	}
	itemAreas = nil
}

func (g *Grid) addGFX(pp *graphics.Pipeline, area image.Rectangle, z int) {
	if g.Skipped() {
		return
	}

	area, _ = g.itemImpl.adjustInnerArea(area)
	g.itemImpl.addGFX(pp, area, z)
	// Decorations
	for _, d := range g.itemImpl.decorations {
		if d.Visible() {
			// Note: mark them as z+1 so that they're drawn on top
			// of the grid's content
			d.addGFX(pp, area, z+1)
		}
	}
	// Compute children items areas
	unitX, unitY := area.Dx()/g.width, area.Dy()/g.height
	//remX := float64(inner.Dx()%g.width) / float64(g.width)
	//remY := float64(inner.Dy()%g.height) / float64(g.height)
	itemAreas := map[Item]image.Rectangle{}
	for y, row := range g.grid {
		for x, i := range row {
			// Skip empty cells
			if i == nil {
				continue
			}
			// Process cell's item area
			if _, ok := itemAreas[i]; !ok {
				itemAreas[i] = image.Rectangle{
					Min: image.Pt(math.MaxInt32, math.MaxInt32),
					Max: image.Pt(math.MinInt32, math.MinInt32),
				}
			}
			// Padding offsets
			//xoff := int(float64(x) * remX)
			//yoff := int(float64(y) * remY)
			r := itemAreas[i]
			r.Min.X = min(r.Min.X, area.Min.X+x*unitX)
			r.Min.Y = min(r.Min.Y, area.Min.Y+y*unitY)
			r.Max.X = max(r.Max.X, area.Min.X+(x+1)*unitX)
			r.Max.Y = max(r.Max.Y, area.Min.Y+(y+1)*unitY)
			itemAreas[i] = r
		}
	}
	// Draw children base items and their own children
	for i, r := range itemAreas {
		i.addGFX(pp, r, z+1)
	}
	itemAreas = nil
}

func (g *Grid) GridBase() *Grid {
	return g
}

func (g *Grid) Columns() int {
	return g.width
}

func (g *Grid) Rows() int {
	return g.height
}

func (g *Grid) Dimensions() (int, int) {
	return g.width, g.height
}

func (g *Grid) Add(x, y, columns, rows int, item Item) {
	if x < 0 || y < 0 || x+columns > g.width || y+rows > g.height {
		panic("ui: cannot add an item outside of the grid")
	}
	if columns == 0 || rows == 0 {
		panic("ui: columns and rows must be > 0")
	}
	// Hack: need to set the address once
	item.setAddr(item)

	g.items = append(g.items, item)
	deletions := map[Item]struct{}{}
	// Add the item
	for ty := y; ty < y+rows; ty++ {
		for tx := x; tx < x+columns; tx++ {
			if idel := g.grid[ty][tx]; idel != nil {
				deletions[idel] = struct{}{}
			}
			g.grid[ty][tx] = item
		}
	}
	// Delete items which space got compromised
	for idel := range deletions {
		for i, curr := range g.items {
			if curr == idel {
				g.items[i] = nil
				g.items = slices.Delete(g.items, i, i+1)
			}
		}
		for ty := range g.grid {
			for tx, curr := range g.grid[ty] {
				if curr == idel {
					g.grid[ty][tx] = nil
				}
			}
		}
	}
	deletions = nil
}

func (g *Grid) ForEach(fn func(i Item)) {
	for _, item := range g.items {
		fn(item)
	}
}

// Options

type GridOption = func(*Grid)

func (g *Grid) WithOptions(opts ...GridOption) *Grid {
	for _, o := range opts {
		o(g)
	}
	return g
}
