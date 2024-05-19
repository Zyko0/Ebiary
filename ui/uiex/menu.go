package uiex

type MenuBar struct {
	*grid
}

func NewMenuBar(columns int) *MenuBar {
	mb := &MenuBar{}
	mb.grid = newGrid(columns, 1, mb)
	return mb
}

// Options

type MenuBarOption func(*MenuBar)

func (mb *MenuBar) WithOptions(opts ...MenuBarOption) *MenuBar {
	for _, o := range opts {
		o(mb)
	}
	return mb
}
