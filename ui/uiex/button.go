package uiex

type Button struct {
	*block
}

func NewButton() *Button {
	b := &Button{}
	b.block = newBlock(b)
	defaultTheme.apply(b)

	return b
}

// Options

type ButtonOption func(*Button)

func (b *Button) WithOptions(opts ...ButtonOption) *Button {
	for _, o := range opts {
		o(b)
	}
	return b
}
