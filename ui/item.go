package ui

import (
	"image"
	"image/color"

	"github.com/Zyko0/Ebiary/ui/internal/core"
	"github.com/Zyko0/Ebiary/ui/internal/graphics"
)

type Item interface {
	// Internal
	base() *itemImpl
	setAddr(item Item)
	update(c *context, area image.Rectangle, z int)
	addGFX(pp *graphics.Pipeline, area image.Rectangle, z int)

	SetItemOptions(opts ...ItemOption)

	// State
	Pressed() bool
	PressedTicks() uint
	JustPressed() bool
	Hovered() bool
	HoveredTicks() uint
	JustUnhovered() bool
	Focused() bool
	JustFocused() bool

	Skipped() bool
	SetSkipped(skipped bool)

	// Style options
	SetShape(shape Shape)
	SetRounding(factor float64)
	SetMinWidth(width int)
	SetMinHeight(height int)
	SetMaxWidth(width int)
	SetMaxHeight(height int)
	SetAlign(x, y AlignMode)
	SetAlignX(x AlignMode)
	SetAlignY(y AlignMode)
	SetMargin(left, right, top, bottom float64)
	SetMarginLeft(pixels float64)
	SetMarginRight(pixels float64)
	SetMarginTop(pixels float64)
	SetMarginBottom(pixels float64)
	SetPadding(left, right, top, bottom int)
	SetPaddingLeft(pixels int)
	SetPaddingRight(pixels int)
	SetPaddingTop(pixels int)
	SetPaddingBottom(pixels int)
	SetBorderWidth(width float64)
	SetBorderColor(clr color.Color)
	SetBorderColorAlpha(alpha float64)
	SetColorPrimary(clr color.Color)
	SetColorSecondary(clr color.Color)
	SetColorPrimaryOffset(offset float64)
	SetColorAlpha(alpha float64)
	SetColorFilling(filling ColorFilling)
	// Getters
	Shape() Shape
	Rounding() float64
	MinSize() (int, int)
	MaxSize() (int, int)
	Align() (AlignMode, AlignMode)
	Margin() (left, right, top, bottom float64)
	Padding() (left, right, top, bottom int)
	Border() (float64, color.Color)
	PrimaryColor() color.Color
	SecondaryColor() color.Color
	PrimaryColorOffset() float64
	Alpha() float64
	ColorFilling() ColorFilling

	// Event
	SetFocusHandling(handled bool)
	SetEventHandling(handled bool)
	SetEventHandler(handler EventHandler)
	EventHandler() EventHandler
	SetEventStyleOptions(EventOptions)
	SetEventActionOptions(EventOptions)
	DoEvent(event Event)

	// Class attributes
	Classes() []string
	SetClasses(classes ...string)
	AddClasses(classes ...string)

	// Decorations
	Decorations() []Decoration
	AddDeco(deco Decoration)

	// User data
	Data() any
	SetData(data any)
}

type Shape graphics.Shape

const (
	ShapeBox     = Shape(graphics.ShapeBox)
	ShapeEllipse = Shape(graphics.ShapeEllipse)
	ShapeRhombus = Shape(graphics.ShapeRhombus)
	ShapeOctogon = Shape(graphics.ShapeOctogon)
	ShapeNone    = Shape(graphics.ShapeNone)
)

type AlignMode byte

const (
	AlignCenter AlignMode = iota
	AlignMin
	AlignMax
	AlignOffset
)

type geom struct {
	lastCursor     image.Point
	lastRegion     image.Rectangle
	lastFullRegion image.Rectangle

	srcOffset     image.Point
	minSize       image.Point
	maxSize       image.Point
	alignX        AlignMode
	alignY        AlignMode
	marginLeft    float64
	marginRight   float64
	marginTop     float64
	marginBottom  float64
	paddingLeft   int
	paddingRight  int
	paddingTop    int
	paddingBottom int
}

type stateFuncs struct {
	defaultFunc      func()
	hoverFunc        func()
	unhoverFunc      func()
	pressFunc        func()
	pressHoverFunc   func()
	releaseFunc      func()
	releaseHoverFunc func()
}

func (sf *stateFuncs) do(st Event) {
	var fn func()
	switch st {
	case Hover:
		fn = sf.hoverFunc
	case Unhover:
		fn = sf.unhoverFunc
	case Press:
		fn = sf.pressFunc
	case PressHover:
		fn = sf.pressHoverFunc
	case Release:
		fn = sf.releaseFunc
	case ReleaseHover:
		fn = sf.releaseHoverFunc
	default:
		fn = sf.defaultFunc
	}
	if fn == nil {
		fn = sf.defaultFunc
	}
	if fn != nil {
		fn()
	}
}

func (sf *stateFuncs) set(ii *itemImpl, opts EventOptions) {
	if fn := opts[Default]; fn != nil {
		sf.defaultFunc = func() { fn(ii.addr) }
	}
	if fn := opts[Hover]; fn != nil {
		sf.hoverFunc = func() { fn(ii.addr) }
	}
	if fn := opts[Unhover]; fn != nil {
		sf.unhoverFunc = func() { fn(ii.addr) }
	}
	if fn := opts[Press]; fn != nil {
		sf.pressFunc = func() { fn(ii.addr) }
	}
	if fn := opts[PressHover]; fn != nil {
		sf.pressHoverFunc = func() { fn(ii.addr) }
	}
	if fn := opts[Release]; fn != nil {
		sf.releaseFunc = func() { fn(ii.addr) }
	}
	if fn := opts[ReleaseHover]; fn != nil {
		sf.releaseHoverFunc = func() { fn(ii.addr) }
	}
}

type state struct {
	noFocus bool
	noEvent bool
	skipped bool

	handler EventHandler

	style  stateFuncs
	action stateFuncs

	styleFunc  func(InputState)
	actionFunc func(InputState)

	justUnhovered bool
	hoveredTicks  uint
	pressedTicks  uint
	focused       bool
	justFocused   bool
}

type body struct {
	shape          graphics.Shape
	alpha          float32
	colorMin       color.Color
	colorMax       color.Color
	colorMinOffset float32
	colorFilling   graphics.ColorFilling
	rounding       float32
}

type border struct {
	alpha float32
	color color.Color
	width float32
}

type itemImpl struct {
	core.ClassesImpl

	addr Item

	body        body
	geom        geom
	state       state
	border      border
	decorations []Decoration

	userData any
}

func newItem(addr Item) *itemImpl {
	return &itemImpl{
		addr: addr,
		body: body{
			alpha: 1.,
		},
	}
}

func (ii *itemImpl) base() *itemImpl {
	return ii
}

func (ii *itemImpl) setAddr(i Item) {
	ii.addr = core.Addr(i)
	for _, d := range ii.decorations {
		if dp, ok := d.(decoWithParent); ok {
			dp.SetParent(i)
		}
	}
}

func (ii *itemImpl) EventHandler() EventHandler {
	if ii.state.handler != nil {
		return ii.state.handler
	}

	return GetEventHandler()
}

func (ii *itemImpl) doStyle(is InputState) {
	if ii.state.styleFunc != nil {
		ii.state.styleFunc(is)
		return
	}
	st := ii.EventHandler().State(ii.addr, is)
	ii.state.style.do(st)
}

func (ii *itemImpl) doAction(is InputState) {
	if ii.state.actionFunc != nil {
		ii.state.actionFunc(is)
		return
	}
	st := ii.EventHandler().State(ii.addr, is)
	ii.state.action.do(st)
}

func (ii *itemImpl) recordState(is InputState) {
	st := ii.EventHandler().State(ii.addr, is)
	switch {
	case st == ReleaseHover, st == Release:
		ii.state.pressedTicks = 0
	case ii.state.pressedTicks > 0:
		ii.state.pressedTicks++
	case st == PressHover, st == Press:
		ii.state.pressedTicks = 1
	}
	// Focus
	if ii.state.pressedTicks == 1 && !ii.state.focused && !ii.state.noFocus {
		SetFocusedItem(ii.addr)
	}
}

func (ii *itemImpl) DoEvent(event Event) {
	ii.state.style.do(event)
	ii.state.action.do(event)
}

func (ii *itemImpl) adjustAreaSize(area image.Rectangle) image.Rectangle {
	if x := ii.geom.minSize.X; x > 0 {
		area.Max.X = max(area.Max.X, area.Min.X+x)
	}
	if y := ii.geom.minSize.Y; y > 0 {
		area.Max.Y = max(area.Max.Y, area.Min.Y+y)
	}
	if x := ii.geom.maxSize.X; x > 0 {
		area.Max.X = min(area.Max.X, area.Min.X+x)
	}
	if y := ii.geom.maxSize.Y; y > 0 {
		area.Max.Y = min(area.Max.Y, area.Min.Y+y)
	}

	return area
}

func (ii *itemImpl) adjustAreaPadding(area image.Rectangle) image.Rectangle {
	area.Min.X += ii.geom.paddingLeft
	area.Min.Y += ii.geom.paddingTop
	area.Max.X -= ii.geom.paddingRight
	area.Max.Y -= ii.geom.paddingBottom

	return area
}

func (ii *itemImpl) alignOffset(area, inner image.Rectangle) (float64, float64) {
	var x, y float64
	switch ii.geom.alignX {
	case AlignMin, AlignOffset:
		x = 0
	case AlignMax:
		x = float64(area.Max.X - inner.Dx() - area.Min.X)
	case AlignCenter:
		x = float64(area.Dx())/2 - float64(inner.Dx())/2
	}
	switch ii.geom.alignY {
	case AlignMin, AlignOffset:
		y = 0
	case AlignMax:
		y = float64(area.Max.Y - inner.Dy() - area.Min.Y)
	case AlignCenter:
		y = float64(area.Dy()/2) - float64(inner.Dy())/2
	}

	return x, y
}

func (ii *itemImpl) adjustInnerArea(area image.Rectangle) (image.Rectangle, image.Rectangle) {
	area = area.Add(ii.geom.srcOffset)
	inner := ii.adjustAreaSize(area)
	x, y := ii.alignOffset(area, inner)
	inner = inner.Add(image.Pt(int(x), int(y)))
	inner = ii.adjustAreaPadding(inner)
	clamped := image.Rect(
		min(max(area.Min.X, inner.Min.X), area.Max.X),
		min(max(area.Min.Y, inner.Min.Y), area.Max.Y),
		max(min(area.Max.X, inner.Max.X), area.Min.X),
		max(min(area.Max.Y, inner.Max.Y), area.Min.Y),
	)

	return clamped, inner
}

func (ii *itemImpl) update(c *context, item Item, area image.Rectangle, z int) {
	// Adjust area with margin, padding
	area.Min.X += int(ii.geom.marginLeft)
	area.Min.Y += int(ii.geom.marginTop)
	area.Max.X -= int(ii.geom.marginRight)
	area.Max.Y -= int(ii.geom.marginBottom)
	// Skip if not event handling
	if ii.state.noEvent {
		return
	}
	ii.state.justUnhovered = false
	ii.state.justFocused = false
	if z >= c.Z && c.Cursor.In(area) {
		if c.Hovered != nil && c.Hovered.Hovered() {
			c.Unhovered = c.Hovered
		}
		c.Hovered = item
		c.Z = z
	}
	if item.Hovered() && c.Hovered != item {
		c.Unhovered = item
	}

	c.DeferredUpdates = append(c.DeferredUpdates, ii.doStyle, ii.doAction, ii.recordState)
}

func (ii *itemImpl) addGFX(pp *graphics.Pipeline, area image.Rectangle, z int) {
	if ii.body.shape == graphics.ShapeNone {
		pp.EnsureLayers(z)
		return
	}

	pp.Add(&graphics.ItemPrimitive{
		Z: z,

		Shape:          ii.body.shape,
		ColorMin:       ii.body.colorMin,
		ColorMax:       ii.body.colorMax,
		ColorMinFactor: ii.body.colorMinOffset,
		ColorFilling:   ii.body.colorFilling,
		ColorAlpha:     ii.body.alpha,
		Rounding:       ii.body.rounding,
		BorderColor:    ii.border.color,
		BorderWidth:    ii.border.width,

		MarginLeft:   float32(ii.geom.marginLeft),
		MarginRight:  float32(ii.geom.marginRight),
		MarginTop:    float32(ii.geom.marginTop),
		MarginBottom: float32(ii.geom.marginBottom),
	}, area)
}

func (ii *itemImpl) LastCursor() image.Point {
	return ii.geom.lastCursor
}

func (ii *itemImpl) LastRegion() image.Rectangle {
	return ii.geom.lastRegion
}

func (ii *itemImpl) LastFullRegion() image.Rectangle {
	return ii.geom.lastFullRegion
}

// State

func (ii *itemImpl) Skipped() bool {
	return ii.state.skipped
}

func (ii *itemImpl) SetSkipped(skipped bool) {
	ii.state.skipped = skipped
}

func (ii *itemImpl) Hovered() bool {
	return ii.state.hoveredTicks > 0
}

func (ii *itemImpl) HoveredTicks() uint {
	return ii.state.hoveredTicks
}

func (ii *itemImpl) JustUnhovered() bool {
	return ii.state.justUnhovered
}

func (ii *itemImpl) Pressed() bool {
	return ii.state.pressedTicks > 0
}

func (ii *itemImpl) PressedTicks() uint {
	return ii.state.pressedTicks
}

func (ii *itemImpl) JustPressed() bool {
	return ii.state.pressedTicks == 1
}

func (ii *itemImpl) Focused() bool {
	return ii.state.focused
}

func (ii *itemImpl) JustFocused() bool {
	return ii.state.justFocused
}

// Options

func (ii *itemImpl) SetItemOptions(opts ...ItemOption) {
	for _, o := range opts {
		o(ii.addr)
	}
}

type ItemOption func(Item)

// Style options

func (ii *itemImpl) SetEventStyleOptions(opts EventOptions) {
	ii.state.style.set(ii, opts)
}

func (ii *itemImpl) SetShape(shape Shape) {
	ii.body.shape = graphics.Shape(shape)
}

func (ii *itemImpl) SetRounding(factor float64) {
	ii.body.rounding = float32(factor)
}

func (ii *itemImpl) SetMinWidth(width int) {
	ii.geom.minSize.X = width
}

func (ii *itemImpl) SetMinHeight(height int) {
	ii.geom.minSize.Y = height
}

func (ii *itemImpl) SetMaxWidth(width int) {
	ii.geom.maxSize.X = width
}

func (ii *itemImpl) SetMaxHeight(height int) {
	ii.geom.maxSize.Y = height
}

func (ii *itemImpl) SetAlign(x, y AlignMode) {
	ii.geom.alignX = x
	ii.geom.alignY = y
}

func (ii *itemImpl) SetAlignX(x AlignMode) {
	ii.geom.alignX = x
}

func (ii *itemImpl) SetAlignY(y AlignMode) {
	ii.geom.alignY = y
}

func (ii *itemImpl) SetMargin(left, right, top, bottom float64) {
	ii.geom.marginLeft = left
	ii.geom.marginRight = right
	ii.geom.marginTop = top
	ii.geom.marginBottom = bottom
}

func (ii *itemImpl) SetMarginLeft(pixels float64) {
	ii.geom.marginLeft = pixels
}

func (ii *itemImpl) SetMarginRight(pixels float64) {
	ii.geom.marginRight = pixels
}

func (ii *itemImpl) SetMarginTop(pixels float64) {
	ii.geom.marginTop = pixels
}

func (ii *itemImpl) SetMarginBottom(pixels float64) {
	ii.geom.marginBottom = pixels
}

func (ii *itemImpl) SetPadding(left, right, top, bottom int) {
	ii.geom.paddingLeft = left
	ii.geom.paddingRight = right
	ii.geom.paddingTop = top
	ii.geom.paddingBottom = bottom
}

func (ii *itemImpl) SetPaddingLeft(pixels int) {
	ii.geom.paddingLeft = pixels
}

func (ii *itemImpl) SetPaddingRight(pixels int) {
	ii.geom.paddingRight = pixels
}

func (ii *itemImpl) SetPaddingTop(pixels int) {
	ii.geom.paddingTop = pixels
}

func (ii *itemImpl) SetPaddingBottom(pixels int) {
	ii.geom.paddingBottom = pixels
}

func (ii *itemImpl) SetBorderWidth(width float64) {
	ii.border.width = float32(width)
}

func (ii *itemImpl) SetBorderColor(clr color.Color) {
	ii.border.color = clr
}

func (ii *itemImpl) SetBorderColorAlpha(alpha float64) {
	ii.border.alpha = float32(alpha)
}

func (ii *itemImpl) SetColorPrimary(clr color.Color) {
	ii.body.colorMin = clr
}

func (ii *itemImpl) SetColorSecondary(clr color.Color) {
	ii.body.colorMax = clr
}

func (ii *itemImpl) SetColorPrimaryOffset(offset float64) {
	ii.body.colorMinOffset = float32(offset)
}

func (ii *itemImpl) SetColorAlpha(alpha float64) {
	ii.body.alpha = float32(alpha)
}

type ColorFilling graphics.ColorFilling

const (
	ColorFillingVertical = ColorFilling(graphics.ColorFillingVertical)
	ColorFillingDistance = ColorFilling(graphics.ColorFillingDistance)
	ColorFillingNone     = ColorFilling(graphics.ColorFillingNone)
)

func (ii *itemImpl) SetColorFilling(filling ColorFilling) {
	ii.body.colorFilling = graphics.ColorFilling(filling)
}

// Getters

func (ii *itemImpl) Shape() Shape {
	return Shape(ii.body.shape)
}

func (ii *itemImpl) Rounding() float64 {
	return float64(ii.body.rounding)
}

func (ii *itemImpl) MinSize() (int, int) {
	return ii.geom.minSize.X, ii.geom.minSize.Y
}

func (ii *itemImpl) MaxSize() (int, int) {
	return ii.geom.maxSize.X, ii.geom.maxSize.Y
}

func (ii *itemImpl) Align() (AlignMode, AlignMode) {
	return AlignMode(ii.geom.alignX), AlignMode(ii.geom.alignY)
}

func (ii *itemImpl) Margin() (left, right, top, bottom float64) {
	return ii.geom.marginLeft,
		ii.geom.marginRight,
		ii.geom.marginTop,
		ii.geom.marginBottom
}

func (ii *itemImpl) Padding() (left, right, top, bottom int) {
	return ii.geom.paddingLeft,
		ii.geom.paddingRight,
		ii.geom.paddingTop,
		ii.geom.paddingBottom
}

func (ii *itemImpl) Border() (float64, color.Color) {
	return float64(ii.border.width), ii.border.color
}

func (ii *itemImpl) PrimaryColor() color.Color {
	return ii.body.colorMin
}

func (ii *itemImpl) SecondaryColor() color.Color {
	return ii.body.colorMax
}

func (ii *itemImpl) PrimaryColorOffset() float64 {
	return float64(ii.body.colorMinOffset)
}

func (ii *itemImpl) Alpha() float64 {
	return float64(ii.body.alpha)
}

func (ii *itemImpl) ColorFilling() ColorFilling {
	return ColorFilling(ii.body.colorFilling)
}

// Event

// SetFocusHandling defines whether an item should take the focused
// state, therefore unfocusing the previously focused item.
func (ii *itemImpl) SetFocusHandling(handled bool) {
	ii.state.noFocus = !handled
}

func (ii *itemImpl) SetEventHandling(handled bool) {
	ii.state.noEvent = !handled
}

func (ii *itemImpl) SetEventActionOptions(opts EventOptions) {
	ii.state.action.set(ii, opts)
}

func (ii *itemImpl) SetEventHandler(handler EventHandler) {
	ii.state.handler = handler
}

// User data

func (ii *itemImpl) Data() any {
	return ii.userData
}

func (ii *itemImpl) SetData(data any) {
	ii.userData = data
}

// Item decorations

func (ii *itemImpl) AddDeco(deco Decoration) {
	// Hack: need to set the address once
	deco.setAddr(deco)

	ii.decorations = append(ii.decorations, deco)
}

func (ii *itemImpl) Decorations() []Decoration {
	return ii.decorations
}

func (ii *itemImpl) SourceOffset() image.Point {
	return ii.geom.srcOffset
}

func (ii *itemImpl) SetSourceOffset(offset image.Point) {
	ii.geom.srcOffset = offset
}

// Base item

func WithData(data any) ItemOption {
	return func(i Item) {
		i.SetData(data)
	}
}

// WithCustomStyleFunc specifies a function to define the style of an item
// based on events.
// The function fn will replace the default function.
// This function is called every tick, and will always called before the
// update function.
func WithCustomStyleFunc[T Item](fn func(T, InputState)) ItemOption {
	return func(i Item) {
		i.base().state.styleFunc = func(is InputState) {
			fn(i.base().addr.(T), is)
		}
	}
}

// WithCustomUpdateFunc specifies a function to be called every tick for logic
// handling.
// The function fn will replace the default function.
// This function can also be used to edit the style of components and will
// overwrite any changes made by the func provided by WithStyleFunc, since it
// executes after.
func WithCustomUpdateFunc[T Item](fn func(T, InputState)) ItemOption {
	return func(i Item) {
		i.base().state.actionFunc = func(is InputState) {
			fn(i.base().addr.(T), is)
		}
	}
}
