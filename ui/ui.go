package ui

import (
	"bytes"
	"image"
	"io/fs"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/font/gofont/goregular"
)

type InputState interface {
	Cursor() image.Point
	MouseWheel() (float64, float64)

	KeyPressDuration(ebiten.Key) int
	MouseButtonPressDuration(ebiten.MouseButton) int
	GamepadButtonPressDuration(ebiten.GamepadID, ebiten.GamepadButton) int
	StandardGamepadButtonPressDuration(ebiten.GamepadID, ebiten.StandardGamepadButton) int
	TouchPressDuration(ebiten.TouchID) int

	KeyJustReleased(ebiten.Key) bool
	MouseButtonJustReleased(ebiten.MouseButton) bool
	GamepadButtonJustReleased(ebiten.GamepadID, ebiten.GamepadButton) bool
	StandardGamepadButtonJustReleased(ebiten.GamepadID, ebiten.StandardGamepadButton) bool
	TouchJustReleased(ebiten.TouchID) bool

	DroppedFiles() fs.FS
}

type inputStateImpl struct {
}

func (inputStateImpl) Cursor() image.Point {
	x, y := ebiten.CursorPosition()
	return image.Point{x, y}
}

func (inputStateImpl) MouseWheel() (float64, float64) {
	return ebiten.Wheel()
}

func (inputStateImpl) KeyPressDuration(key ebiten.Key) int {
	return inpututil.KeyPressDuration(key)
}

func (inputStateImpl) MouseButtonPressDuration(button ebiten.MouseButton) int {
	return inpututil.MouseButtonPressDuration(button)
}

func (inputStateImpl) GamepadButtonPressDuration(id ebiten.GamepadID, button ebiten.GamepadButton) int {
	return inpututil.GamepadButtonPressDuration(id, button)
}

func (inputStateImpl) StandardGamepadButtonPressDuration(id ebiten.GamepadID, button ebiten.StandardGamepadButton) int {
	return inpututil.StandardGamepadButtonPressDuration(id, button)
}

func (inputStateImpl) TouchPressDuration(id ebiten.TouchID) int {
	return inpututil.TouchPressDuration(id)
}

func (inputStateImpl) KeyJustReleased(key ebiten.Key) bool {
	return inpututil.IsKeyJustReleased(key)
}

func (inputStateImpl) MouseButtonJustReleased(button ebiten.MouseButton) bool {
	return inpututil.IsMouseButtonJustReleased(button)
}

func (inputStateImpl) GamepadButtonJustReleased(id ebiten.GamepadID, button ebiten.GamepadButton) bool {
	return inpututil.IsGamepadButtonJustReleased(id, button)
}

func (inputStateImpl) StandardGamepadButtonJustReleased(id ebiten.GamepadID, button ebiten.StandardGamepadButton) bool {
	return inpututil.IsStandardGamepadButtonJustReleased(id, button)
}

func (inputStateImpl) TouchJustReleased(id ebiten.TouchID) bool {
	return inpututil.IsTouchJustReleased(id)
}

func (inputStateImpl) DroppedFiles() fs.FS {
	return ebiten.DroppedFiles()
}

// Global state

var (
	uiInputState   inputStateImpl = inputStateImpl{}
	uiEventHandler EventHandler   = eventHandlerImpl{}

	uiFS = os.DirFS(".")

	uiFontSource *text.GoTextFaceSource
	uiFontSize   = 14.

	uiBlockOption func(*Block)
	uiGridOption  func(*Grid)

	uiFocusedItem Item
)

func init() {
	var err error

	uiFontSource, err = text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		log.Fatalf("ui: cannot load goregular default font: %v", err)
	}
}

func GetInputState() InputState {
	return uiInputState
}

func GetEventHandler() EventHandler {
	return uiEventHandler
}

func SetEventHandler(ec EventHandler) {
	uiEventHandler = ec
}

func GetFS() fs.FS {
	return uiFS
}

func SetFS(fsys fs.FS) {
	uiFS = fsys
}

func GetFontSource() *text.GoTextFaceSource {
	return uiFontSource
}

func SetFontSource(source *text.GoTextFaceSource) {
	uiFontSource = source
}

func GetFontSize() float64 {
	return uiFontSize
}

func SetFontSize(size float64) {
	uiFontSize = size
}

func SetBlockOption(option func(*Block)) {
	uiBlockOption = option
}

func SetGridOption(option func(*Grid)) {
	uiGridOption = option
}

func FocusedItem() Item {
	return uiFocusedItem
}

func SetFocusedItem(item Item) {
	// Unfocus previously focused item
	if uiFocusedItem != nil {
		uiFocusedItem.base().state.justFocused = false
		uiFocusedItem.base().state.focused = false
	}
	item.base().state.justFocused = true
	item.base().state.focused = true
	uiFocusedItem = item.base().addr
}
