package opt

import (
	"github.com/Zyko0/Ebiary/ui"
	"github.com/Zyko0/Ebiary/ui/uiex"
)

type (
	pictureImage struct{}

	picture struct {
		Image pictureImage
	}
)

var Picture picture

func (picture) Options(opts ...ui.ItemOption) uiex.PictureOption {
	return func(p *uiex.Picture) {
		p.SetItemOptions(opts...)
	}
}

func (pictureImage) Options(opts ...uiex.ImageOption) uiex.PictureOption {
	return func(p *uiex.Picture) {
		p.Image().WithOptions(opts...)
	}
}
