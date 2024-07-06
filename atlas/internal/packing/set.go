package packing

import (
	"image"
	"slices"
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

		minSize: image.Pt(1, 1),
	}
	if opts != nil {
		s.minSize = opts.MinSize
	}

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
	s.empties = s.empties[:0]
	for i := range s.tmps {
		if s.tmps[i].Dx() < s.minSize.X || s.tmps[i].Dy() < s.minSize.Y {
			continue
		}
		var contained bool
		// Filter out any duplicate or any empty region that is already
		// contained by another one
		for _, empty := range s.empties {
			if s.tmps[i] == empty || s.tmps[i].In(empty) {
				contained = true
				break
			}
		}
		if contained {
			continue
		}
		s.empties = append(s.empties, s.tmps[i])
	}

	//fmt.Println("len:", len(s.rects))

	return true
}

func (s *Set) Free(rect *image.Rectangle) {
	if rect == nil || len(s.rects) == 0 {
		return
	}
	idx := slices.Index(s.rects, rect)
	if idx != -1 {
		s.rects = slices.Delete(s.rects, idx, idx+1)

		s.tmps = s.tmps[:0]
		for _, e := range s.empties {
			if e.Max.X == rect.Min.X || e.Min.X == rect.Max.X || e.Max.Y == rect.Min.Y || e.Min.Y == rect.Max.Y {
				s.tmps = append(s.tmps, e)
			}
		}
		println("count around:", len(s.tmps))
		// Create a big rectangle containing all neighbours
		parent := *rect
		for _, e := range s.tmps {
			parent.Min.X = min(parent.Min.X, e.Min.X)
			parent.Min.Y = min(parent.Min.Y, e.Min.Y)
			parent.Max.X = max(parent.Max.X, e.Max.X)
			parent.Max.Y = max(parent.Max.Y, e.Max.Y)
		}
		//s.tmps = append(s.tmps, parent)
		var occupied []image.Rectangle
		for _, r := range s.rects {
			if r.In(parent) {
				occupied = append(occupied, *r)
			}
		}
		// Merge a maximum of rectangles around the freed one
		/*var done bool
		for !done {
			done = true
			biggest := *rect
			toDel := image.Rectangle{}
			bs := biggest.Dx() * biggest.Dy()
			for _, e := range s.tmps {
				ok := true
				t := *rect
				if e.Max.X > rect.Max.X && e.Min.Y <= rect.Min.Y && e.Max.Y >= rect.Max.Y {
					t.Max.X = e.Max.X
					for _, r := range occupied {
						if r.In(t) {
							ok = false
							break
						}
					}
				} else if e.Max.X < rect.Min.X && e.Min.Y <= rect.Min.Y && e.Max.Y >= rect.Max.Y {
					t.Min.X = e.Min.X
					for _, r := range occupied {
						if r.In(t) {
							ok = false
							break
						}
					}
				}
				if e.Max.Y > rect.Min.Y && e.Min.X <= rect.Min.X && e.Max.X >= rect.Max.X {
					t.Min.Y = e.Min.Y
					for _, r := range occupied {
						if r.In(t) {
							ok = false
							break
						}
					}
				} else if e.Max.Y < rect.Min.Y && e.Min.X <= rect.Min.X && e.Max.X >= rect.Max.X {
					t.Min.X = e.Min.X
					for _, r := range occupied {
						if r.In(t) {
							ok = false
							break
						}
					}
				}
				if ok {
					if size := t.Dx() * t.Dy(); size > bs {
						bs = size
						biggest = t
						toDel = e
						done = false
					}
				}
			}
			if biggest != *rect {
				idx := slices.Index(s.empties, toDel)
				if idx != -1 {
					s.empties = slices.Delete(s.empties, idx, idx+1)
				}
				//s.empties = append(s.empties, biggest)
				s.tmps = append(s.tmps[:0], biggest)
				done = true
			} else {
				s.tmps = append(s.tmps[:0], *rect)
				done = true
			}
		}*/
		// Prepare the empty regions for next insertion
		for i := range s.tmps {
			if s.tmps[i].Dx() < s.minSize.X || s.tmps[i].Dy() < s.minSize.Y {
				continue
			}
			var contained bool
			// Filter out any duplicate or any empty region that is already
			// contained by another one
			for _, e := range s.empties {
				if s.tmps[i] == e || s.tmps[i].In(e) {
					contained = true
					break
				}
			}
			if contained {
				continue
			}
			s.empties = append(s.empties, s.tmps[i])
		}
		//s.empties = append(s.empties, s.tmps...)
	}
}
