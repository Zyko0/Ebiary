package opt

import (
	"github.com/Zyko0/Ebiary/ui"
	"github.com/Zyko0/Ebiary/ui/uiex"
)

type (
	caretBar struct{}

	caret struct{}

	textInput struct {
		Caret caret
		textItem[*uiex.TextInput]
	}
)

var TextInput textInput

func (textInput) Options(opts ...ui.ItemOption) uiex.TextInputOption {
	return func(ti *uiex.TextInput) {
		ti.SetItemOptions(opts...)
	}
}

func (textInput) Multiline(multiline bool) uiex.TextInputOption {
	return func(ti *uiex.TextInput) {
		ti.SetMultiline(multiline)
	}
}

func (caret) Options(opts ...ui.ItemOption) uiex.TextInputOption {
	return func(ti *uiex.TextInput) {
		ti.SetCaretOptions(opts...)
	}
}

func (caret) Thickness(thickness int) uiex.TextInputOption {
	return func(ti *uiex.TextInput) {
		ti.SetCaretThickness(thickness)
	}
}
