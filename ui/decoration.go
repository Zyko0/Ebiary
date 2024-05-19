package ui

type Decoration interface {
	Item

	Visible() bool
}

type decoWithParent interface {
	Decoration

	SetParent(parent Item)
}
