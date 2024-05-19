package ui

import "github.com/hajimehoshi/ebiten/v2"

type Event byte

const (
	Default Event = iota
	Hover
	Unhover
	Press
	PressHover
	JustPress
	Release
	ReleaseHover
	EventMax
)

type EventOptions [EventMax]ItemOption

type EventHandler interface {
	State(item Item, is InputState) Event
}

type eventHandlerImpl struct{}

func (ec eventHandlerImpl) State(item Item, is InputState) Event {
	click := is.MouseButtonPressDuration(ebiten.MouseButtonLeft)
	switch {
	case item.Hovered():
		switch {
		case item.Pressed() && click >= 1:
			return PressHover
		case !item.Pressed() && click == 1:
			return PressHover
		case item.Pressed() && click == 0:
			return ReleaseHover
		case click == 0:
			return Hover
		default:
			return Default
		}
	case item.Pressed() && click == 0:
		return Release
	case item.JustUnhovered():
		return Unhover
	case item.Pressed():
		return Press
	default:
		return Default
	}
}
