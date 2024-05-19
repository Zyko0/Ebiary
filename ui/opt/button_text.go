package opt

import (
	"github.com/Zyko0/Ebiary/ui"
	"github.com/Zyko0/Ebiary/ui/uiex"
)

type (
	buttonText struct {
		textItem[*uiex.ButtonText]
	}
)

var ButtonText buttonText

func (buttonText) Options(opts ...ui.ItemOption) uiex.ButtonTextOption {
	return func(bt *uiex.ButtonText) {
		bt.SetItemOptions(opts...)
	}
}
