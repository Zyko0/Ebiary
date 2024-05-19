package uiex

import (
	"github.com/Zyko0/Ebiary/ui"
)

// Item

type StyleItem[T ui.Item] []ui.ItemOption

func (s StyleItem[T]) Apply(item T) {
	if s == nil {
		return
	}
	for _, o := range s {
		o(item)
	}
}

func (s *StyleItem[T]) Option() func(T) {
	return s.Apply
}

// Content

type StyleContent[T ui.Content] []func(T)

func (s StyleContent[T]) Apply(content T) {
	if s == nil {
		return
	}
	for _, o := range s {
		o(content)
	}
}

func (s *StyleContent[T]) Option() func(T) {
	return s.Apply
}

// Theme

type Theme struct {
	// Base items
	Block StyleItem[*ui.Block]
	Grid  StyleItem[*ui.Grid]
	// Items
	Button     StyleItem[*Button]
	ButtonText StyleItem[*ButtonText]
	Label      StyleItem[*Label]
	Bar        StyleItem[*Bar]
	TextInput  StyleItem[*TextInput]
	// Decorations
	Scrollbar StyleItem[*Scrollbar]
	// Content
	BasicText StyleContent[*BasicText]
	RichText  StyleContent[*RichText]
	Image     StyleContent[*Image]
}

type (
	uiBlock     interface{ BlockBase() *ui.Block }
	uiGrid      interface{ GridBase() *ui.Grid }
	withContent interface{ Content() ui.Content }
)

func (t *Theme) apply(item ui.Item) {
	// Apply the base styles first
	switch ii := item.(type) {
	case uiBlock:
		t.Block.Apply(ii.BlockBase())
	case uiGrid:
		t.Grid.Apply(ii.GridBase())
	default:
		return
	}
	// Item implementations
	switch ii := item.(type) {
	case *Button:
		t.Button.Apply(ii)
	case *ButtonText:
		t.ButtonText.Apply(ii)
	case *Label:
		t.Label.Apply(ii)
	case *Bar:
		t.Bar.Apply(ii)
	case *TextInput:
		t.TextInput.Apply(ii)
	// Decorations
	case *Scrollbar:
		t.Scrollbar.Apply(ii)
	}
	// Content implementations
	if ic, ok := item.(withContent); ok {
		switch c := ic.Content().(type) {
		case *BasicText:
			t.BasicText.Apply(c)
		case *RichText:
			t.RichText.Apply(c)
		case *Image:
			t.Image.Apply(c)
		}
	}
}

func (t *Theme) Apply(layout *ui.Layout) {
	if t == nil {
		return
	}
	layout.Grid().ForEach(t.apply)
}

type (
	ClassStyleItem[T ui.Item]       map[string]StyleItem[T]
	ClassStyleContent[T ui.Content] map[string]StyleContent[T]
)

func (csi ClassStyleItem[T]) apply(i T, classes []string) {
	for _, class := range classes {
		if s, ok := csi[class]; ok {
			s.Apply(i)
		}
	}
}

func (csi ClassStyleContent[T]) apply(c T, classes []string) {
	for _, class := range classes {
		if s, ok := csi[class]; ok {
			s.Apply(c)
		}
	}
}

type ClassTheme struct {
	// Base theme
	*Theme
	// Base items class styles
	Block ClassStyleItem[*ui.Block]
	Grid  ClassStyleItem[*ui.Grid]
	// Items class styles
	Button     ClassStyleItem[*Button]
	ButtonText ClassStyleItem[*ButtonText]
	Label      ClassStyleItem[*Label]
	Bar        ClassStyleItem[*Bar]
	TextInput  ClassStyleItem[*TextInput]
	// Decorations
	Scrollbar ClassStyleItem[*Scrollbar]
	// Content class styles
	BasicText ClassStyleContent[*BasicText]
	RichText  ClassStyleContent[*RichText]
	Image     ClassStyleContent[*Image]
}

func (t *ClassTheme) apply(item ui.Item) {
	// Apply base theme
	t.Theme.apply(item)
	// Apply class styles
	classes := item.Classes()
	// Apply the base styles first
	switch ii := item.(type) {
	case uiBlock:
		t.Block.apply(ii.BlockBase(), classes)
	case uiGrid:
		t.Grid.apply(ii.GridBase(), classes)
		ii.GridBase().ForEach(t.apply)
	default:
		return
	}

	// Item implementations
	switch ii := item.(type) {
	case *Button:
		t.Button.apply(ii, classes)
	case *ButtonText:
		t.ButtonText.apply(ii, classes)
	case *Label:
		t.Label.apply(ii, classes)
	case *Bar:
		t.Bar.apply(ii, classes)
	case *TextInput:
		t.TextInput.apply(ii, classes)
	// Decorations
	case *Scrollbar:
		t.Scrollbar.apply(ii, classes)
	}
	// Content implementations
	if ic, ok := item.(withContent); ok {
		switch c := ic.Content().(type) {
		case *BasicText:
			t.BasicText.apply(c, c.Classes())
		case *RichText:
			t.RichText.apply(c, c.Classes())
		case *Image:
			t.Image.apply(c, c.Classes())
		}
	}
}

func (t *ClassTheme) Apply(layout *ui.Layout) {
	if t == nil {
		return
	}

	t.apply(layout.Grid())
}

// Default

var (
	defaultTheme *ClassTheme
)

func SetTheme(t *Theme) {
	defaultTheme = &ClassTheme{
		Theme: t,
	}
	ui.SetBlockOption(t.Block.Apply)
	ui.SetGridOption(t.Grid.Apply)
}

func SetClassTheme(t *ClassTheme) {
	defaultTheme = t
	ui.SetBlockOption(func(b *ui.Block) {
		t.Theme.Block.Apply(b)
		t.Block.apply(b, b.Classes())
	})
	ui.SetGridOption(func(g *ui.Grid) {
		t.Theme.Grid.Apply(g)
		t.Grid.apply(g, g.Classes())
	})
}

func init() {
	defaultTheme = &ClassTheme{
		Theme: &Theme{},
	}
	/*defaultTheme = &ClassTheme{
		Theme: &Theme{
			Block: []ui.ItemOption{
				func(i ui.Item) {
					i.SetColorPrimary(color.RGBA{32, 32, 32, 255})
					i.SetColorSecondary(color.RGBA{32, 32, 32, 255})
					i.SetBorderWidth(1)
					i.SetBorderColor(color.RGBA{128, 128, 128, 255})
				},
			},
			Grid:       []ui.ItemOption{},
			Button:     []ui.ItemOption{},
			ButtonText: []ui.ItemOption{},
			Label:      []ui.ItemOption{},
			Bar:        []ui.ItemOption{},
			// Decorations
			Scrollbar: []ui.ItemOption{
				func(i ui.Item) {
					sb := i.(*Scrollbar)
					sb.SetColorPrimary(color.RGBA{32, 32, 32, 255})
					sb.SetColorSecondary(color.RGBA{32, 32, 32, 255})
					sb.SetColorFilling(ui.ColorFillingNone)
					sb.SetBorderWidth(1)
					sb.SetBorderColor(color.Black)
					sb.WithWidth(24)
					sb.WithDirection(DirectionVertical)
					sb.SetEventStyleOptions(ui.EventOptions{
						ui.Unhover: func(_ ui.Item) {
							alpha := max(sb.cursor.Alpha()-5./255, 0)
							sb.cursor.SetColorAlpha(alpha)
							sb.cursor.SetColorPrimary(color.RGBA{32, 32, 32, 255})
							sb.cursor.SetColorSecondary(color.RGBA{32, 32, 32, 255})
						},
						ui.ReleaseHover: func(_ ui.Item) {
							alpha := max(sb.cursor.Alpha()-5./255, 0)
							sb.cursor.SetColorAlpha(alpha)
							sb.cursor.SetColorPrimary(color.RGBA{32, 32, 32, 255})
							sb.cursor.SetColorSecondary(color.RGBA{32, 32, 32, 255})
						},
						ui.Hover: func(_ ui.Item) {
							alpha := min(sb.cursor.Alpha()+10./255, 0.5)
							sb.cursor.SetColorAlpha(alpha)
						},
					})
					sb.WithCursorBarOptions(
						func(cursor *Bar) {
							cursor.SetEventStyleOptions(ui.EventOptions{
								ui.Unhover: func(_ ui.Item) {
									cursor.SetColorPrimary(color.RGBA{32, 32, 32, 255})
									cursor.SetColorSecondary(color.RGBA{32, 32, 32, 255})
								},
								ui.ReleaseHover: func(_ ui.Item) {
									cursor.SetColorPrimary(color.RGBA{32, 32, 32, 255})
									cursor.SetColorSecondary(color.RGBA{32, 32, 32, 255})
								},
								ui.Hover: func(_ ui.Item) {
									cursor.SetColorPrimary(color.RGBA{64, 64, 64, 255})
									cursor.SetColorSecondary(color.RGBA{64, 64, 64, 255})
									cursor.SetColorAlpha(0.5)
								},
								ui.Press: func(_ ui.Item) {
									cursor.SetColorPrimary(color.RGBA{96, 96, 96, 255})
									cursor.SetColorSecondary(color.RGBA{96, 96, 96, 255})
									cursor.SetColorAlpha(0.5)
								},
							})
							cursor.SetColorPrimary(color.RGBA{64, 64, 64, 255})
							cursor.SetColorSecondary(color.RGBA{64, 64, 64, 255})
							cursor.SetColorAlpha(0)
							cursor.WithDirection(DirectionVertical)
						},
					)
				},
			},
			// Content
			Text:  []func(*Text){},
			Image: []func(*Image){},
		},
	}*/
}
