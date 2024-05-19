package uiex

import (
	"image"

	"github.com/Zyko0/Ebiary/ui"
)

type Scrollbar struct {
	*Bar

	parent Scrollable
	cursor *Bar

	drag    image.Point
	dragged bool

	alwaysVisible bool
	wheelSpeed    float64
}

func NewScrollbar() *Scrollbar {
	sb := &Scrollbar{
		Bar: NewBar(),

		cursor: NewBar(),

		alwaysVisible: false,
	}
	// Hack: need to set the addr to the scrollbar type since it's composed by *Bar
	sb.block.addr = sb
	sb.SetFocusHandling(false)
	sb.cursor.SetAlign(ui.AlignOffset, ui.AlignOffset)
	sb.cursor.SetFocusHandling(false)
	sb.AddDeco(sb.cursor)
	defaultTheme.apply(sb)

	// Scrollbar moving logic (when not hovering cursor)
	sb.SetItemOptions(
		ui.WithCustomUpdateFunc(func(isb *Scrollbar, is ui.InputState) {
			// Calculate cursor length
			global := isb.parent.LastFullRegion()
			local := isb.parent.LastRegion()
			cs := isb.LastCursor()
			var length float64
			if isb.direction == DirectionHorizontal {
				cs.Y = 0
				length = min(float64(local.Dx())/float64(global.Dx()), 1)
				length *= float64(isb.LastRegion().Dx())
				isb.cursor.SetAlign(ui.AlignOffset, ui.AlignCenter)
			} else {
				cs.X = 0
				length = min(float64(local.Dy())/float64(global.Dy()), 1)
				length *= float64(isb.LastRegion().Dy())
				sb.cursor.SetAlign(ui.AlignCenter, ui.AlignOffset)
			}
			isb.cursor.SetLength(int(length))
			// Offset logic
			poffset := isb.parent.SourceOffset()
			offset := isb.cursor.SourceOffset().Add(cs.Sub(sb.drag))
			// If not dragged, recompute the offset based on parent regions
			if !sb.dragged {
				offset = poffset
				if sb.direction == DirectionHorizontal {
					offset.Y = 0
					offset.X = int(0.5 + float64(-offset.X*sb.LastRegion().Dx())/float64(global.Dx()))
				} else {
					offset.X = 0
					offset.Y = int(0.5 + float64(-offset.Y*sb.LastRegion().Dy())/float64(global.Dy()))
				}
			}
			// If dragging logic is already handled, skip
			var barClick bool
			if !isb.dragged && isb.Pressed() {
				offset = isb.cursor.SourceOffset()
				if isb.direction == DirectionHorizontal {
					x := cs.X - int(length)
					offset.X = min(max(x, 0), isb.LastRegion().Dx()-int(length))
				} else {
					y := cs.Y - int(length)
					offset.Y = min(max(y, 0), isb.LastRegion().Dy()-int(length))
				}
				barClick = true
			}
			// Wheel logic
			var wheel bool
			var wheelX, wheelY float64
			if isb.parent.Focused() {
				wheelX, wheelY = is.MouseWheel()
				wheelX, wheelY = max(min(wheelX, 1), -1), max(min(-wheelY, 1), -1)
			}
			dx, dy := 10., 10.
			// Ajust cursor and parent's source offsets
			if sb.direction == DirectionHorizontal {
				wheel = wheelX != 0
				offset = offset.Add(image.Pt(int(wheelX*dx), 0))
				offset.X = min(max(offset.X, 0), sb.LastRegion().Dx()-sb.cursor.length)
				px := float64(offset.X) / float64(sb.LastRegion().Dx()-sb.cursor.length)
				px *= float64(sb.parent.LastFullRegion().Dx() - sb.parent.LastRegion().Dx())
				poffset.X = -int(px)
			} else {
				wheel = wheelY != 0
				offset = offset.Add(image.Pt(0, int(wheelY*dy)))
				offset.Y = min(max(offset.Y, 0), sb.LastRegion().Dy()-sb.cursor.length)
				py := float64(offset.Y) / float64(sb.LastRegion().Dy()-sb.cursor.length)
				py *= float64(sb.parent.LastFullRegion().Dy() - sb.parent.LastRegion().Dy())
				poffset.Y = -int(py)
			}
			// If the new cursor position didn't affect the offset, don't retain
			// new dragging position
			// TODO: seems useless
			/*if offset == sb.cursor.SourceOffset() {
				return
			}*/
			isb.cursor.SetSourceOffset(offset)
			if sb.dragged || wheel || barClick {
				isb.parent.SetSourceOffset(poffset)
			}
		}),
	)
	// Cursor bar logic (hovered)
	sb.cursor.SetEventActionOptions(ui.EventOptions{
		ui.PressHover: func(_ ui.Item) {
			cs := sb.LastCursor()
			offset := sb.cursor.SourceOffset().Add(cs.Sub(sb.drag))
			if sb.direction == DirectionHorizontal {
				cs.Y, offset.Y = 0, 0
				offset.X = min(max(offset.X, 0), sb.LastRegion().Dx()-sb.cursor.length)
			} else {
				cs.X, offset.X = 0, 0
				offset.Y = min(max(offset.Y, 0), sb.LastRegion().Dy()-sb.cursor.length)
			}
			// Do not retain new drag position if it didn't affect cursor offset
			if !sb.dragged || offset != sb.cursor.SourceOffset() {
				sb.drag = cs
			}
			sb.dragged = true
		},
		ui.Press: func(_ ui.Item) {
			sb.cursor.DoEvent(ui.PressHover)
		},
		ui.Release: func(_ ui.Item) {
			sb.dragged = false
		},
		ui.ReleaseHover: func(_ ui.Item) {
			sb.dragged = false
		},
	})

	return sb
}

func (sb *Scrollbar) WithClasses(classes ...string) *Scrollbar {
	sb.SetClasses(classes...)
	return sb
}

func (sb *Scrollbar) Visible() bool {
	if sb.alwaysVisible {
		return true
	}

	offset := sb.parent.SourceOffset()
	var visible bool
	if sb.direction == DirectionHorizontal {
		visible = sb.parent.LastFullRegion().Dx() != sb.parent.LastRegion().Dx()
		offset.X = 0
	} else {
		visible = sb.parent.LastFullRegion().Dy() != sb.parent.LastRegion().Dy()
		offset.Y = 0
	}
	// If parent does not need scrollbar, scroll the content back to the origin
	if !visible {
		sb.parent.SetSourceOffset(offset)
	}
	return visible
}

func (sb *Scrollbar) SetParent(parent ui.Item) {
	sb.parent = parent.(Scrollable)
}

func (sb *Scrollbar) SetVisiblity(visible bool) *Scrollbar {
	sb.visible = visible
	return sb
}

func (sb *Scrollbar) SetDirection(direction Direction) *Scrollbar {
	sb.Bar.SetDirection(direction)
	switch direction {
	case DirectionHorizontal:
		sb.Bar.SetAlignY(ui.AlignMax)
	case DirectionVertical:
		sb.Bar.SetAlignX(ui.AlignMax)
	}

	sb.cursor.SetDirection(direction)

	return sb
}

// Options

type ScrollbarOption func(*Scrollbar)

func (sb *Scrollbar) WithOptions(opts ...ScrollbarOption) *Scrollbar {
	for _, o := range opts {
		o(sb)
	}
	return sb
}

func (sb *Scrollbar) WithBarOptions(opts ...BarOption) *Scrollbar {
	sb.Bar.WithOptions(opts...)
	return sb
}

func (sb *Scrollbar) WithCursorOptions(opts ...ui.ItemOption) *Scrollbar {
	sb.cursor.SetItemOptions(opts...)
	return sb
}

func (sb *Scrollbar) WithCursorBarOptions(opts ...BarOption) *Scrollbar {
	sb.cursor.WithOptions(opts...)
	return sb
}
