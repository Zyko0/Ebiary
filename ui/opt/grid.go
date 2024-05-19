package opt

import (
	"github.com/Zyko0/Ebiary/ui"
)

type grid struct{}

var Grid grid

func (grid) Options(opts ...ui.ItemOption) ui.GridOption {
	return func(b *ui.Grid) {
		b.SetItemOptions(opts...)
	}
}
