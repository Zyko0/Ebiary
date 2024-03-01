package asset

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/bmp"
)

func newShader[T any](data []byte) (T, error) {
	s, err := ebiten.NewShader(data)

	return any(s).(T), err
}

func newImagePNG[T any](data []byte) (T, error) {
	var t T
	i, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return t, err
	}
	img := ebiten.NewImageFromImage(i)

	return any(img).(T), nil
}

func newImageJPG[T any](data []byte) (T, error) {
	var t T
	i, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		return t, err
	}
	img := ebiten.NewImageFromImage(i)

	return any(img).(T), nil
}

func newImageGIF[T any](data []byte) (T, error) {
	var t T
	i, err := gif.Decode(bytes.NewReader(data))
	if err != nil {
		return t, err
	}
	img := ebiten.NewImageFromImage(i)

	return any(img).(T), nil
}

func newImageBMP[T any](data []byte) (T, error) {
	var t T
	i, err := bmp.Decode(bytes.NewReader(data))
	if err != nil {
		return t, err
	}
	img := ebiten.NewImageFromImage(i)

	return any(img).(T), nil
}

// Loaders

type Loader func([]byte) (any, error)

func NewImage(data []byte) (any, error) {
	img, _, err := image.Decode(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	return ebiten.NewImageFromImage(img), nil
}

func NewImageWithOptions(opts *ebiten.NewImageFromImageOptions) Loader {
	return func(data []byte) (any, error) {
		img, _, err := image.Decode(bytes.NewBuffer(data))
		if err != nil {
			return nil, err
		}
		return ebiten.NewImageFromImageWithOptions(img, opts), nil
	}
}

func NewShader(data []byte) (any, error) {
	return ebiten.NewShader(data)
}
