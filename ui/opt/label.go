package opt

import (
	"github.com/Zyko0/Ebiary/ui"
	"github.com/Zyko0/Ebiary/ui/uiex"
)

type (
	label struct {
		textItem[*uiex.Label]
	}
)

var Label label

func (label) Options(opts ...ui.ItemOption) uiex.LabelOption {
	return func(l *uiex.Label) {
		l.SetItemOptions(opts...)
	}
}
