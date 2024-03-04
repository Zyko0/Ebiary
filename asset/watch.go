package asset

import (
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
)

var (
	onceWatch      sync.Once
	watcherInitErr error
	watcher        *fsnotify.Watcher

	trackedFiles *sync.Map
)

// Assets watcher

type asset interface {
	lock()
	unlock()
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

func watchAsset(path string, obj asset) error {
	onceWatch.Do(watcherInit)
	if watcherInitErr != nil {
		return watcherInitErr
	}
	// Store reference to the watched asset
	path = filepath.Clean(path)
	files, _ := trackedFiles.LoadOrStore(path, &sync.Map{})
	files.(*sync.Map).Store(obj, struct{}{})
	// Read file content and ensure the path leads to a file
	f, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	// Perform the initial asset loading
	err = obj.refresh(f, nil)
	if err != nil {
		return err
	}
	// Add to watcher
	err = watcher.Add(path)
	if err != nil {
		return err
	}

	return nil
}

func unwatchAsset(path string, obj asset) error {
	onceWatch.Do(watcherInit)
	if watcherInitErr != nil {
		return watcherInitErr
	}
	// Lock object to prevent Value() access in case of LiveAsset
	obj.lock()
	defer obj.unlock()
	// Remove from tracked files
	path = filepath.Clean(path)
	files, _ := trackedFiles.Load(path)
	filesMap := files.(*sync.Map)
	filesMap.Delete(obj)
	// Unwatch the path if no more tracked files for it
	count := 0
	filesMap.Range(func(key, value any) bool {
		count++
		return true
	})
	if count == 0 {
		err := watcher.Remove(path)
		if err != nil && !errors.Is(err, fsnotify.ErrNonExistentWatch) {
			return err
		}
		trackedFiles.Delete(path)
	}

	return nil
}
