package packing

import (
	"fmt"
	"image"
)

type Set struct {
	width  int
	height int
	rects  []*image.Rectangle
	frees  []image.Rectangle
	tmps   []image.Rectangle

	minSize image.Point
}

type NewSetOptions struct {
	MinSize image.Point
}

func NewSet(width, height int, opts *NewSetOptions) *Set {
	s := &Set{
		width:  width,
		height: height,
		frees: []image.Rectangle{
			image.Rect(0, 0, width, height),
		},

		minSize: image.Pt(1, 1),
	}
	if opts != nil {
		s.minSize = opts.MinSize
	}

	return s
}

func appendFreeRects(rects []image.Rectangle, parent, filled image.Rectangle) []image.Rectangle {
	if !filled.In(parent) {
		return append(rects, parent)
	}
	if filled.Min.X > parent.Min.X {
		rects = append(rects, image.Rect(
			parent.Min.X,
			parent.Min.Y,
			filled.Min.X,
			parent.Max.Y,
		))
	}
	if filled.Max.X < parent.Max.X {
		rects = append(rects, image.Rect(
			filled.Max.X,
			parent.Min.Y,
			parent.Max.X,
			parent.Max.Y,
		))
	}
	if filled.Min.Y > parent.Min.Y {
		rects = append(rects, image.Rect(
			parent.Min.X,
			parent.Min.Y,
			parent.Max.X,
			filled.Min.Y,
		))
	}
	if filled.Max.Y < parent.Max.Y {
		rects = append(rects, image.Rect(
			parent.Min.X,
			filled.Max.Y,
			parent.Max.X,
			parent.Max.Y,
		))
	}
	return rects
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func (s *Set) Insert(rect *image.Rectangle) bool {
	s.tmps = append(s.tmps[:0], s.frees...)
	// Filter too small rectangles
	n := 0
	for i := 0; i < len(s.tmps); i++ {
		if rect.Dx() > s.tmps[i].Dx() || rect.Dy() > s.tmps[i].Dy() {
			continue
		}
		s.tmps[n] = s.tmps[i]
		n++
	}
	s.tmps = s.tmps[:n]
	// TODO: find the most optimal in free ones
	if len(s.tmps) == 0 {
		return false
	}
	// Find best rect
	best := s.tmps[0]
	bs := best.Min.X + best.Min.Y
	for i := range s.tmps {
		if d := s.tmps[i].Min.X + s.tmps[i].Min.Y; d < bs {
			best = s.tmps[i]
			bs = d
		}
	}

	slot := best
	*rect = rect.Add(slot.Min)
	s.rects = append(s.rects, rect)

	s.tmps = s.tmps[:0]
	for i := range s.frees {
		if ix := rect.Intersect(s.frees[i]); !ix.Empty() {
			s.tmps = appendFreeRects(s.tmps, s.frees[i], ix)
		} else {
			s.tmps = append(s.tmps, s.frees[i])
		}
	}

	s.frees = s.frees[:0]
	for i := range s.tmps {
		if s.tmps[i].Dx() < s.minSize.X || s.tmps[i].Dy() < s.minSize.Y {
			continue
		}
		var contained bool
		for _, parent := range s.frees {
			if s.tmps[i] == parent || s.tmps[i].In(parent) {
				contained = true
				break
			}
		}
		if contained {
			continue
		}
		s.frees = append(s.frees, s.tmps[i])
	}

	fmt.Println("len:", len(s.rects))

	return true
}
