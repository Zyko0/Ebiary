package packing

import (
	"image"
	"slices"
	"sort"
)

type Set struct {
	width   int
	height  int
	rects   []*image.Rectangle
	empties []image.Rectangle
	tmps    []image.Rectangle

	minSize image.Point
}

type NewSetOptions struct {
	MinSize image.Point
}

func NewSet(width, height int, opts *NewSetOptions) *Set {
	s := &Set{
		width:  width,
		height: height,
		empties: []image.Rectangle{
			image.Rect(0, 0, width, height),
		},
	}
	if opts != nil {
		s.minSize = opts.MinSize
	}
	s.minSize.X = max(s.minSize.X, 1)
	s.minSize.Y = max(s.minSize.Y, 1)

	return s
}

func appendEmptyNeighbours(rects []image.Rectangle, parent, filled image.Rectangle) []image.Rectangle {
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

func (s *Set) sanitizeEmptyRegions(current []image.Rectangle) {
	s.empties = s.empties[:0]
	// Sort empty regions by size to ensure smaller regions are evicted if
	// contained by bigger ones
	sort.SliceStable(current, func(i, j int) bool {
		si := current[i].Dx() * current[i].Dy()
		sj := current[j].Dx() * current[j].Dy()
		return si > sj
	})
	for i := range current {
		if current[i].Dx() < s.minSize.X || current[i].Dy() < s.minSize.Y {
			continue
		}
		var contained bool
		// Filter out any duplicate or any empty region that is already
		// contained by another one
		for _, empty := range s.empties {
			if current[i] == empty || current[i].In(empty) {
				contained = true
				break
			}
		}
		if contained {
			continue
		}
		s.empties = append(s.empties, current[i])
	}
}

func (s *Set) Insert(rect *image.Rectangle) bool {
	// Set the free regions from last insertion
	s.tmps = append(s.tmps[:0], s.empties...)
	// Filter out too small regions
	n := 0
	for i := 0; i < len(s.tmps); i++ {
		if rect.Dx() > s.tmps[i].Dx() || rect.Dy() > s.tmps[i].Dy() {
			continue
		}
		s.tmps[n] = s.tmps[i]
		n++
	}
	s.tmps = s.tmps[:n]
	// Abort if no available rectangle
	if len(s.tmps) == 0 {
		return false
	}
	// Find best rectangle (the closest to top left corner)
	best := s.tmps[0]
	bs := best.Min.X + best.Min.Y
	for i := range s.tmps {
		if d := s.tmps[i].Min.X + s.tmps[i].Min.Y; d < bs {
			best = s.tmps[i]
			bs = d
		}
	}
	// Insert the provided rectangle with the origin of the best free region
	*rect = rect.Add(best.Min)
	s.rects = append(s.rects, rect)
	// Split the regions that used to be empty with new empty neighbours
	s.tmps = s.tmps[:0]
	for i := range s.empties {
		if ix := rect.Intersect(s.empties[i]); !ix.Empty() {
			s.tmps = appendEmptyNeighbours(s.tmps, s.empties[i], ix)
		} else {
			s.tmps = append(s.tmps, s.empties[i])
		}
	}
	// Prepare the empty regions for next insertion
	s.sanitizeEmptyRegions(s.tmps)

	return true
}

func (s *Set) Free(rect *image.Rectangle) {
	if rect == nil || len(s.rects) == 0 {
		return
	}
	idx := slices.Index(s.rects, rect)
	if idx != -1 {
		s.rects = slices.Delete(s.rects, idx, idx+1)

		// Try to grow the just freed region until it's not possible anymore
		// TODO: awful algorithm but curiously fast enough for the moment
		// Filter Y-intersecting rectangles only
		freed := *rect
		s.tmps = s.tmps[:0]
		for _, r := range s.rects {
			if r.Max.Y >= freed.Min.Y && r.Min.Y <= freed.Max.Y {
				s.tmps = append(s.tmps, *r)
			}
		}
		for freed.Min.X > 0 {
			freed.Min.X -= 1
			var found bool
			for _, r := range s.tmps {
				if !freed.Intersect(r).Empty() {
					found = true
					break
				}
			}
			if found {
				freed.Min.X += 1
				break
			}
		}
		for freed.Max.X < s.width {
			freed.Max.X += 1
			var found bool
			for _, r := range s.tmps {
				if !freed.Intersect(r).Empty() {
					found = true
					break
				}
			}
			if found {
				freed.Max.X -= 1
				break
			}
		}
		// Filter X-intersecting rectangles only
		s.tmps = s.tmps[:0]
		for _, r := range s.rects {
			if r.Max.X >= freed.Min.X && r.Min.X <= freed.Max.X {
				s.tmps = append(s.tmps, *r)
			}
		}
		for freed.Min.Y > 0 {
			freed.Min.Y -= 1
			var found bool
			for _, r := range s.tmps {
				if !freed.Intersect(r).Empty() {
					found = true
					break
				}
			}
			if found {
				freed.Min.Y += 1
				break
			}
		}
		for freed.Max.Y < s.height {
			freed.Max.Y += 1
			var found bool
			for _, r := range s.tmps {
				if !freed.Intersect(r).Empty() {
					found = true
					break
				}
			}
			if found {
				freed.Max.Y -= 1
				break
			}
		}
		f0 := freed

		// SECOND

		freed = *rect
		// Filter X-intersecting rectangles only
		s.tmps = s.tmps[:0]
		for _, r := range s.rects {
			if r.Max.X >= freed.Min.X && r.Min.X <= freed.Max.X {
				s.tmps = append(s.tmps, *r)
			}
		}
		for freed.Min.Y > 0 {
			freed.Min.Y -= 1
			var found bool
			for _, r := range s.tmps {
				if !freed.Intersect(r).Empty() {
					found = true
					break
				}
			}
			if found {
				freed.Min.Y += 1
				break
			}
		}
		for freed.Max.Y < s.height {
			freed.Max.Y += 1
			var found bool
			for _, r := range s.tmps {
				if !freed.Intersect(r).Empty() {
					found = true
					break
				}
			}
			if found {
				freed.Max.Y -= 1
				break
			}
		}
		s.tmps = s.tmps[:0]
		for _, r := range s.rects {
			if r.Max.Y >= freed.Min.Y && r.Min.Y <= freed.Max.Y {
				s.tmps = append(s.tmps, *r)
			}
		}
		for freed.Min.X > 0 {
			freed.Min.X -= 1
			var found bool
			for _, r := range s.tmps {
				if !freed.Intersect(r).Empty() {
					found = true
					break
				}
			}
			if found {
				freed.Min.X += 1
				break
			}
		}
		for freed.Max.X < s.width {
			freed.Max.X += 1
			var found bool
			for _, r := range s.tmps {
				if !freed.Intersect(r).Empty() {
					found = true
					break
				}
			}
			if found {
				freed.Max.X -= 1
				break
			}
		}
		f1 := freed

		s0 := f0.Dx() * f0.Dy()
		s1 := f1.Dx() * f1.Dy()
		//fmt.Println("s0", s0, f0.Dx(), f0.Dy(), "s1", s1, f1.Dx(), f1.Dy())
		if s0 > s1 {
			freed = f0
		} else {
			freed = f1
		}

		s.empties = append(s.empties, freed)
		s.tmps = append(s.tmps[:0], s.empties...)

		// Sanitize empty space
		s.sanitizeEmptyRegions(s.tmps)
	}
}
