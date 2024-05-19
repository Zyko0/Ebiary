package uitext

import (
	"image/color"
	"time"
)

type TextMask byte

const (
	TextMaskDefault       TextMask = 0
	TextMaskBold          TextMask = 1
	TextMaskItalic        TextMask = 2
	TextMaskUnderline     TextMask = 4
	TextMaskStrikethrough TextMask = 8
	TextMaskBackground    TextMask = 16
)

type TextEffect struct {
	Start   int
	End     int
	Mask    TextMask
	ColorFg color.Color
	ColorBg color.Color
}

type Stack struct {
	effects []*TextEffect
	index   int
}

func NewEffectStack() *Stack {
	return &Stack{
		effects: []*TextEffect{
			{
				End: -1,
			},
		},
		index: 0,
	}
}

func (s *Stack) ResetIndex() {
	s.index = 0
}

var now = time.Now()

func (s *Stack) Effects() []*TextEffect {
	return s.effects
}

func (s *Stack) Effect(index int) *TextEffect {
	var once bool
	for i := s.index; i < len(s.effects); i++ {
		e := s.effects[i]
		if index >= e.Start && (e.End == -1 || index < e.End) {
			s.index = i
			once = true
		} else if once {
			break
		}
	}
	for i := s.index; i > 0; i-- {
		e := s.effects[i]
		if index >= e.Start && (e.End == -1 || index < e.End) {
			s.index = i
			break
		}
	}

	return s.effects[s.index]
}

func (s *Stack) PushFg(index int, clr color.Color) {
	e := *s.effects[s.index]
	e.ColorFg = clr
	e.Start = index
	e.End = -1
	s.effects = append(s.effects, &e)
	s.index = len(s.effects) - 1
}

func (s *Stack) PushBg(index int, clr color.Color) {
	e := *s.effects[s.index]
	e.ColorBg = clr
	e.Mask |= TextMaskBackground
	e.Start = index
	e.End = -1
	s.effects = append(s.effects, &e)
	s.index = len(s.effects) - 1
}

func (s *Stack) PushMask(index int, mask TextMask) {
	e := *s.effects[s.index]
	e.Mask |= mask
	e.Start = index
	e.End = -1
	s.effects = append(s.effects, &e)
	s.index = len(s.effects) - 1
}

func (s *Stack) Pop(index int) {
	for s.index > 0 {
		if s.effects[s.index].End == -1 {
			s.effects[s.index].End = index
			return
		}
		s.index--
	}
}

func (s *Stack) Reset() {
	s.index = 0
	s.effects = s.effects[:1]
}
