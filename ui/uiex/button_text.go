package uiex

type ButtonText Button

func NewButtonText(str string) *ButtonText {
	bt := &ButtonText{}
	bt.block = newBlock(bt)
	bt.block.SetContent(NewBasicText(str))
	defaultTheme.apply(bt)

	return bt
}

func (b *ButtonText) Text() Text {
	return b.Content().(Text)
}

func (b *ButtonText) SetText(txt Text) {
	b.SetContent(txt)
}

// Options

type ButtonTextOption func(*ButtonText)

func (b *ButtonText) WithOptions(opts ...ButtonTextOption) *ButtonText {
	for _, o := range opts {
		o(b)
	}
	return b
}
