package opt

import (
	"github.com/Zyko0/Ebiary/ui"
)

type block struct{}

var Block block

func (block) Options(opts ...ui.ItemOption) ui.BlockOption {
	return func(b *ui.Block) {
		b.SetItemOptions(opts...)
	}
}

func (block) Content(content ui.Content) ui.BlockOption {
	return func(b *ui.Block) {
		b.SetContent(content)
	}
}