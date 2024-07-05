package packing

import (
	"image"
	"slices"
)

type Set struct {
	width  int
	height int
	rects  []*image.Rectangle
	frees  []image.Rectangle

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
	parents := make([]image.Rectangle, len(s.frees))
	copy(parents, s.frees)
	frees := []image.Rectangle{}
	currents := make([]image.Rectangle, len(s.rects))
	for i := range s.rects {
		currents[i] = *s.rects[i]
	}
	var done bool
	var newParents []image.Rectangle
	for !done {
		newParents = newParents[:0]
		done = true
		for _, pr := range parents {
			var contains bool
			//n := 0
			for i, r := range currents {
				if ix := r.Intersect(pr); !ix.Empty() { //r.In(pr) {
					contains = true
					done = false
					newParents = appendFreeRects(newParents, pr, ix)
					if ix == r {
						//currents = slices.Delete(currents, i, i+1)
					}
					//currents = slices.Delete(currents, i, i+1)
					//currents[i] =
					break
				}
				_ = i
				//currents[n] = currents[i]
				//n++
			}
			//currents = currents[:n]
			if pr.Dx() < s.minSize.X || pr.Dy() < s.minSize.Y {
				continue
			}
			if !contains && slices.Index(frees, pr) < 0 {
				frees = append(frees, pr)
			}
		}
		parents = append(parents[:0], newParents...)
	}
	// Filter too small rectangles, as well as duplicates
	n := 0
	//println("len frees before:", len(frees))
	s.frees = append(s.frees[:0], frees...)
	for i := 0; i < len(frees); i++ {
		if rect.Dx() > frees[i].Dx() || rect.Dy() > frees[i].Dy() {
			continue
		}
		frees[n] = frees[i]
		n++
	}
	frees = frees[:n]
	// TODO: find the most optimal in free ones
	_ = frees
	if len(frees) == 0 {
		return false
	}
	// Find best rect
	best := frees[0]
	/*bestDiff := abs(rect.Dx()-best.Dx()) + abs(rect.Dy()-best.Dy())
	for i := 1; i < len(frees); i++ {
		cr := frees[i]
		diff := abs(rect.Dx()-cr.Dx()) + abs(rect.Dy()-cr.Dy())
		if diff < bestDiff {
			best = cr
			bestDiff = diff
		}
	}*/
	//println("number of free rects:", len(frees))
	//fmt.Println("frees:", frees)
	slot := best //frees[0] //frees[rand.Intn(len(frees))]
	*rect = rect.Add(slot.Min)
	s.rects = append(s.rects, rect)
	//println("len s.rects:", len(s.rects))
	// TODO: remove below, should be safe now
	/*for i := range s.rects {
		for j := range s.rects {
			if i != j && s.rects[i].Eq(*s.rects[j]) {
				panic("wtf")
			}
		}
	}*/

	return true
}
