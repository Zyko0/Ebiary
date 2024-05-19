package opt

import (
	"github.com/Zyko0/Ebiary/ui"
	"github.com/Zyko0/Ebiary/ui/uiex"
	"github.com/hajimehoshi/ebiten/v2"
)

type img struct{}

var Image img

// Classes

func (img) Classes(classes ...string) uiex.ImageOption {
	return func(i *uiex.Image) {
		i.SetClasses(classes...)
	}
}

// Image

func (img) FillContainer(fill bool) uiex.ImageOption {
	return func(i *uiex.Image) {
		i.SetContainerFilling(fill)
	}
}

func (img) ColorScale(cs ebiten.ColorScale) uiex.ImageOption {
	return func(i *uiex.Image) {
		i.SetColorScale(cs)
	}
}

// Alignment

func (img) Align(x, y ui.AlignMode) uiex.ImageOption {
	return func(i *uiex.Image) {
		i.SetAlignX(x)
		i.SetAlignY(y)
	}
}

func (img) AlignLeft() uiex.ImageOption {
	return func(i *uiex.Image) {
		i.SetAlignX(ui.AlignMin)
	}
}

func (img) AlignRight() uiex.ImageOption {
	return func(i *uiex.Image) {
		i.SetAlignX(ui.AlignMax)
	}
}

func (img) AlignCenter() uiex.ImageOption {
	return func(i *uiex.Image) {
		i.SetAlignX(ui.AlignCenter)
		i.SetAlignY(ui.AlignCenter)
	}
}

func (img) AlignCenterX() uiex.ImageOption {
	return func(i *uiex.Image) {
		i.SetAlignX(ui.AlignCenter)
	}
}

func (img) AlignTop() uiex.ImageOption {
	return func(i *uiex.Image) {
		i.SetAlignY(ui.AlignMin)
	}
}

func (img) AlignBottom() uiex.ImageOption {
	return func(i *uiex.Image) {
		i.SetAlignY(ui.AlignMax)
	}
}

func (img) AlignCenterY() uiex.ImageOption {
	return func(i *uiex.Image) {
		i.SetAlignY(ui.AlignCenter)
	}
}
