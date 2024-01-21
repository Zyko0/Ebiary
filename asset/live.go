package asset

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/bmp"
)

var (
	watcherInitErr error
	watcher        *fsnotify.Watcher
	trackedFiles   *sync.Map
	onceWatch      sync.Once
)

// Loading functions

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

func newImageBMP[T any](data []byte) (T, error) {
	var t T
	i, err := bmp.Decode(bytes.NewReader(data))
	if err != nil {
		return t, err
	}
	img := ebiten.NewImageFromImage(i)

	return any(img).(T), nil
}

// Assets watcher

type asset interface {
	refresh([]byte, error) error
}

func watcherInit() {
	trackedFiles = &sync.Map{}
	watcher, watcherInitErr = fsnotify.NewWatcher()
	if watcherInitErr != nil {
		return
	}

	// Listen for file events and refresh assets accordingly
	go func() {
		for {
			if event, ok := <-watcher.Events; ok {
				if event.Has(fsnotify.Write) {
					path := filepath.Clean(event.Name)
					files, has := trackedFiles.Load(path)
					if !has {
						continue
					}
					data, err := os.ReadFile(path)
					// Reload each asset registered against the path
					files.(*sync.Map).Range(func(a, _ any) bool {
						a.(asset).refresh(data, err)
						return true
					})
				}
			} else {
				return
			}
		}
	}()
}

type LiveAsset[T any] struct {
	mutex sync.Mutex

	path   string
	err    error
	value  *T
	fnLoad func([]byte) (T, error)
}

func NewLiveAsset[T any](path string) (*LiveAsset[T], error) {
	onceWatch.Do(watcherInit)
	if watcherInitErr != nil {
		return nil, watcherInitErr
	}

	// Register file watching
	a := &LiveAsset[T]{}
	path = filepath.Clean(path)
	files, _ := trackedFiles.LoadOrStore(path, &sync.Map{})
	files.(*sync.Map).Store(a, struct{}{})
	// Read file content and ensure the path leads to a file
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	// Define the loading function
	var v T
	switch any(v).(type) {
	case *ebiten.Shader:
		a.fnLoad = newShader[T]
	case *ebiten.Image:
		ext := filepath.Ext(path)
		switch ext {
		case ".png":
			a.fnLoad = newImagePNG[T]
		case ".jpg", ".jpeg":
			a.fnLoad = newImageJPG[T]
		case ".bmp":
			a.fnLoad = newImageBMP[T]
		default:
			return nil, fmt.Errorf("unknown image extension: '%s'", ext)
		}
	case []byte:
		a.fnLoad = func(b []byte) (T, error) {
			return any(b).(T), nil
		}
	default:
		return nil, fmt.Errorf("unknown asset type: '%s'", reflect.TypeOf(v).String())
	}
	// Perform the initial asset loading
	v, err = a.fnLoad(f)
	if err != nil {
		return nil, err
	}
	a.value = &v
	// Add to watcher
	err = watcher.Add(path)
	if err != nil {
		return nil, err
	}
	// Set a finalizer to unregister the file from tracked file at collection
	runtime.SetFinalizer(a, func(a *LiveAsset[T]) {
		a.Dispose()
	})

	return a, nil
}

func (a *LiveAsset[T]) refresh(data []byte, err error) error {
	if a == nil {
		return nil
	}
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err != nil {
		a.err = err
	} else {
		a.err = nil
		v, err := a.fnLoad(data)
		if err != nil {
			a.err = err
		} else {
			*a.value = v
		}
	}

	return a.err
}

// Value returns the most recent asset value
func (a *LiveAsset[T]) Value() T {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return *a.value
}

// Error returns the last encountered error while reloading the file,
// returns nil if no error
func (a *LiveAsset[T]) Error() error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	return a.err
}

// Dispose unregisters this asset from the file watcher
func (a *LiveAsset[T]) Dispose() {
	if a == nil {
		return
	}
	a.mutex.Lock()
	defer a.mutex.Unlock()

	files, ok := trackedFiles.Load(a.path)
	if ok {
		files.(*sync.Map).Delete(a)
	}
}
