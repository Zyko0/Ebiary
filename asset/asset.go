package asset

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// Loadable defines an interface for assets that can handle their own
// loading / initialization based on the file's content.
// Deserialize should be implemented on any custom type that needs to
// be supported by asset.FS or LiveAsset objects.
type Loadable interface {
	Deserialize([]byte) error
}

// LiveAsset represents a self-reloading asset object, its content is
// accessible by calling the Value() method.
type LiveAsset[T any] struct {
	mutex sync.Mutex

	path   string
	err    error
	obj    T
	fnLoad func([]byte) (T, error)
}

// NewLiveAsset creates an hot-reloadable asset for ebitengine common types.
// The default supported types are *ebiten.Shader and *ebiten.Image (png,bmp,jpg).
func NewLiveAsset[T any](path string) (*LiveAsset[T], error) {
	var fn func([]byte) (T, error)
	var v T
	switch any(v).(type) {
	case *ebiten.Shader:
		fn = newShader[T]
	case *ebiten.Image:
		ext := filepath.Ext(path)
		switch ext {
		case ".png":
			fn = newImagePNG[T]
		case ".jpg", ".jpeg":
			fn = newImageJPG[T]
		case ".gif":
			fn = newImageGIF[T]
		case ".bmp":
			fn = newImageBMP[T]
		default:
			return nil, fmt.Errorf("asset: unknown image extension '%s' for file %s", ext, path)
		}
	case []byte:
		fn = func(b []byte) (T, error) {
			return any(b).(T), nil
		}
	case Loadable:
		fn = func(b []byte) (T, error) {
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
	default:
		return nil, fmt.Errorf("asset: unknown asset type '%s' for file %s", reflect.TypeOf(v).String(), path)
	}

	// Register file watching
	a := &LiveAsset[T]{
		fnLoad: fn,
	}
	// Set a finalizer to unregister the file from tracked file at GC
	runtime.SetFinalizer(a, func(a *LiveAsset[T]) {
		a.Dispose()
	})
	if err := watchAsset(path, a); err != nil {
		return nil, fmt.Errorf("asset: can't watch file %s: %v", path, err)
	}

	return a, nil
}

// NewLiveAssetFunc is a convenience function to create a new LiveAsset
// with the given path and a custom load function.
func NewLiveAssetFunc[T any](path string, fn func([]byte) (T, error)) (*LiveAsset[T], error) {
	// Register file watching
	a := &LiveAsset[T]{
		fnLoad: fn,
	}
	// Set a finalizer to unregister the file from tracked file at collection
	runtime.SetFinalizer(a, func(a *LiveAsset[T]) {
		a.Dispose()
	})
	if err := watchAsset(path, a); err != nil {
		return nil, fmt.Errorf("asset: can't watch file %s: %v", path, err)
	}

	return a, nil
}

func (a *LiveAsset[T]) lock() {
	a.mutex.Lock()
}

func (a *LiveAsset[T]) unlock() {
	a.mutex.Unlock()
}

// refresh reloads the object's content and is called automatically on file update.
func (a *LiveAsset[T]) refresh(data []byte, err error) error {
	if a == nil {
		return nil
	}
	a.lock()
	defer a.unlock()

	if err != nil {
		a.err = err
	} else {
		a.err = nil
		v, err := a.fnLoad(data)
		if err != nil {
			a.err = err
		} else {
			a.obj = v
		}
	}

	return a.err
}

// Value returns the most recent asset value
func (a *LiveAsset[T]) Value() T {
	a.lock()
	defer a.unlock()

	return a.obj
}

// Error returns the last encountered error while reloading the file,
// returns nil if no error
func (a *LiveAsset[T]) Error() error {
	a.lock()
	defer a.unlock()

	return a.err
}

// Dispose unregisters this asset from the file watcher
func (a *LiveAsset[T]) Dispose() {
	if a == nil {
		return
	}

	unwatchAsset(a.path, a)
}
