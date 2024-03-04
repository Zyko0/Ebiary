package asset

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"reflect"

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

func NewShader(data []byte) (any, error) {
	return ebiten.NewShader(data)
}

// Meta loaders

func NewImageWithOptionsLoader(opts *ebiten.NewImageFromImageOptions) Loader {
	return func(data []byte) (any, error) {
		img, _, err := image.Decode(bytes.NewBuffer(data))
		if err != nil {
			return nil, err
		}
		return ebiten.NewImageFromImageWithOptions(img, opts), nil
	}
}

func NewLoadableLoader[T any]() Loader {
	return func(b []byte) (any, error) {
		var obj T
		rt := reflect.TypeOf(obj)
		if rt.Kind() == reflect.Ptr {
			obj = reflect.New(rt.Elem()).Interface().(T)
		}
		err := any(obj).(Loadable).Deserialize(b)
		if err != nil {
			return obj, err
		}
		return obj, nil
	}
}
