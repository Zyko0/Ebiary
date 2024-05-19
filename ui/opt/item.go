package opt

import (
	"image/color"

	"github.com/Zyko0/Ebiary/ui"
)

// Item allows to bundle a list of item implementation-specific options
// as a single ItemOption, for convenience (e.g a theme definition).
// The resulting option is unsafe since it contains implementation-specific
// options that can be passed to the wrong item implementation.
func Item[T ui.Item](opts ...func(T)) ui.ItemOption {
	return func(i ui.Item) {
		ii := i.(T)
		for _, o := range opts {
			o(ii)
		}
	}
}

// Event

func NoEvent() ui.ItemOption {
	return func(i ui.Item) {
		i.SetEventHandling(false)
	}
}

func EventHandler(handler ui.EventHandler) ui.ItemOption {
	return func(i ui.Item) {
		i.SetEventHandler(handler)
	}
}

func EventStyle(opts ui.EventOptions) ui.ItemOption {
	return func(i ui.Item) {
		i.SetEventStyleOptions(opts)
	}
}

func EventAction(opts ui.EventOptions) ui.ItemOption {
	return func(i ui.Item) {
		i.SetEventActionOptions(opts)
	}
}

func DoEvent(event ui.Event) ui.ItemOption {
	return func(i ui.Item) {
		i.DoEvent(event)
	}
}

// Decorations

func Decorations(decorations ...ui.Decoration) ui.ItemOption {
	return func(i ui.Item) {
		for _, d := range decorations {
			i.AddDeco(d)
		}
	}
}

// Classes

func Classes(classes ...string) ui.ItemOption {
	return func(i ui.Item) {
		i.SetClasses(classes...)
	}
}

func Multi(opts ...ui.ItemOption) ui.ItemOption {
	return func(i ui.Item) {
		for _, o := range opts {
			o(i)
		}
	}
}

func Noop() ui.ItemOption {
	return func(_ ui.Item) {}
}

// Shape

func Shape(shape ui.Shape) ui.ItemOption {
	return func(i ui.Item) {
		i.SetShape(shape)
	}
}

func Box() ui.ItemOption {
	return func(i ui.Item) {
		i.SetShape(ui.ShapeBox)
	}
}

func Ellipse() ui.ItemOption {
	return func(i ui.Item) {
		i.SetShape(ui.ShapeEllipse)
	}
}

func Rhombus() ui.ItemOption {
	return func(i ui.Item) {
		i.SetShape(ui.ShapeRhombus)
	}
}

func Octogon() ui.ItemOption {
	return func(i ui.Item) {
		i.SetShape(ui.ShapeOctogon)
	}
}

func NoShape() ui.ItemOption {
	return func(i ui.Item) {
		i.SetShape(ui.ShapeNone)
	}
}

func Rounding(factor float64) ui.ItemOption {
	return func(i ui.Item) {
		i.SetRounding(factor)
	}
}

// Color

func RGB(r, g, b uint8) ui.ItemOption {
	return func(i ui.Item) {
		i.SetColorPrimary(color.RGBA{r, g, b, 255})
		i.SetColorSecondary(color.RGBA{r, g, b, 255})
	}
}

func RGBA(r, g, b, a uint8) ui.ItemOption {
	return func(i ui.Item) {
		i.SetColorPrimary(color.RGBA{r, g, b, 255})
		i.SetColorSecondary(color.RGBA{r, g, b, 255})
		i.SetColorAlpha(float64(a) / 255)
	}
}

func Color(clr color.Color) ui.ItemOption {
	return func(i ui.Item) {
		i.SetColorPrimary(clr)
		i.SetColorSecondary(clr)
	}
}

func Alpha(alpha float64) ui.ItemOption {
	return func(i ui.Item) {
		i.SetColorAlpha(alpha)
	}
}

func AlphaDecay(decr float64, minAlpha float64) ui.ItemOption {
	return func(i ui.Item) {
		alpha := max(i.Alpha()-decr, minAlpha)
		i.SetColorAlpha(alpha)
	}
}

func AlphaIncr(incr float64, maxAlpha float64) ui.ItemOption {
	return func(i ui.Item) {
		alpha := min(i.Alpha()+incr, maxAlpha)
		i.SetColorAlpha(alpha)
	}
}

func Filling(filling ui.ColorFilling) ui.ItemOption {
	return func(i ui.Item) {
		i.SetColorFilling(filling)
	}
}

// Size

func MinWidth(width int) ui.ItemOption {
	return func(i ui.Item) {
		i.SetMinWidth(width)
	}
}

func MinHeight(height int) ui.ItemOption {
	return func(i ui.Item) {
		i.SetMinHeight(height)
	}
}

func MaxWidth(width int) ui.ItemOption {
	return func(i ui.Item) {
		i.SetMaxWidth(width)
	}
}

func MaxHeight(height int) ui.ItemOption {
	return func(i ui.Item) {
		i.SetMaxHeight(height)
	}
}

func MinSize(width, height int) ui.ItemOption {
	return func(i ui.Item) {
		i.SetMinWidth(width)
		i.SetMinHeight(height)
	}
}

func MaxSize(width, height int) ui.ItemOption {
	return func(i ui.Item) {
		i.SetMaxWidth(width)
		i.SetMaxHeight(height)
	}
}

func Size(width, height int) ui.ItemOption {
	return func(i ui.Item) {
		i.SetMinWidth(width)
		i.SetMinHeight(height)
		i.SetMaxWidth(width)
		i.SetMaxHeight(height)
	}
}

func SizeAuto() ui.ItemOption {
	return Size(0, 0)
}

// Alignment

func Align(x, y ui.AlignMode) ui.ItemOption {
	return func(i ui.Item) {
		i.SetAlign(x, y)
	}
}

func AlignLeft() ui.ItemOption {
	return func(i ui.Item) {
		i.SetAlignX(ui.AlignMin)
	}
}

func AlignRight() ui.ItemOption {
	return func(i ui.Item) {
		i.SetAlignX(ui.AlignMax)
	}
}

func AlignCenter() ui.ItemOption {
	return func(i ui.Item) {
		i.SetAlignX(ui.AlignCenter)
		i.SetAlignY(ui.AlignCenter)
	}
}

func AlignCenterX() ui.ItemOption {
	return func(i ui.Item) {
		i.SetAlignX(ui.AlignCenter)
	}
}

func AlignTop() ui.ItemOption {
	return func(i ui.Item) {
		i.SetAlignY(ui.AlignMin)
	}
}

func AlignBottom() ui.ItemOption {
	return func(i ui.Item) {
		i.SetAlignY(ui.AlignMax)
	}
}

func AlignCenterY() ui.ItemOption {
	return func(i ui.Item) {
		i.SetAlignY(ui.AlignCenter)
	}
}

// Margin

func Margin(pixels float64) ui.ItemOption {
	return func(i ui.Item) {
		i.SetMarginLeft(pixels)
		i.SetMarginRight(pixels)
		i.SetMarginTop(pixels)
		i.SetMarginBottom(pixels)
	}
}

func MarginLeft(pixels float64) ui.ItemOption {
	return func(i ui.Item) {
		i.SetMarginLeft(pixels)
	}
}

func MarginRight(pixels float64) ui.ItemOption {
	return func(i ui.Item) {
		i.SetMarginRight(pixels)
	}
}

func MarginTop(pixels float64) ui.ItemOption {
	return func(i ui.Item) {
		i.SetMarginTop(pixels)
	}
}

func MarginBottom(pixels float64) ui.ItemOption {
	return func(i ui.Item) {
		i.SetMarginBottom(pixels)
	}
}

// Padding

func Padding(pixels int) ui.ItemOption {
	return func(i ui.Item) {
		i.SetPaddingLeft(pixels)
		i.SetPaddingRight(pixels)
		i.SetPaddingTop(pixels)
		i.SetPaddingBottom(pixels)
	}
}

func PaddingLeft(pixels int) ui.ItemOption {
	return func(i ui.Item) {
		i.SetPaddingLeft(pixels)
	}
}

func PaddingRight(pixels int) ui.ItemOption {
	return func(i ui.Item) {
		i.SetPaddingRight(pixels)
	}
}

func PaddingTop(pixels int) ui.ItemOption {
	return func(i ui.Item) {
		i.SetPaddingTop(pixels)
	}
}

func PaddingBottom(pixels int) ui.ItemOption {
	return func(i ui.Item) {
		i.SetPaddingBottom(pixels)
	}
}

// Border

func Border(width float64, clr color.Color) ui.ItemOption {
	return func(i ui.Item) {
		i.SetBorderWidth(width)
		i.SetBorderColor(clr)
	}
}

func BorderRGB(r, g, b uint8) ui.ItemOption {
	return func(i ui.Item) {
		i.SetBorderColor(color.RGBA{r, g, b, 255})
	}
}

func BorderRGBA(r, g, b, a uint8) ui.ItemOption {
	return func(i ui.Item) {
		i.SetBorderColor(color.RGBA{r, g, b, 255})
		i.SetBorderColorAlpha(float64(a) / 255)
	}
}

func BorderWidth(width float64) ui.ItemOption {
	return func(i ui.Item) {
		i.SetBorderWidth(width)
	}
}

func BorderColor(clr color.Color) ui.ItemOption {
	return func(i ui.Item) {
		i.SetBorderColor(clr)
	}
}
