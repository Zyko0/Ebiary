package asset

import (
	"embed"
	"fmt"
	"io/fs"
	"runtime"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// Hot reloading assets for a fs.FS

type fsAsset struct {
	mutex sync.Mutex

	obj any
	fn  Loader
}

func (a *fsAsset) lock() {
	a.mutex.Lock()
}

func (a *fsAsset) unlock() {
	a.mutex.Unlock()
}

func (a *fsAsset) refresh(data []byte, err error) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if err != nil {
		return err
	}
	obj, err := a.fn(data)
	if err != nil {
		return err
	}
	a.obj = obj

	return nil
}

// FS

type FS struct {
	fs.FS

	watch   bool
	loaders map[string]Loader
	objects map[string]*fsAsset
}

type LoadRule struct {
	FilePattern string
	Func        Loader
	Overrides   []*LoadRule
}

func initLoadingRules(fsys fs.FS, loaders map[string]Loader, rules []*LoadRule) error {
	if len(rules) == 0 {
		return nil
	}
	for _, rule := range rules {
		// List all files matching pattern
		files, err := fs.Glob(fsys, rule.FilePattern)
		if err != nil {
			return fmt.Errorf("fs: invalid rule pattern '%s': %v", rule.FilePattern, err)
		}
		for _, f := range files {
			loaders[f] = rule.Func
		}
		// Recursive overrides
		if err := initLoadingRules(fsys, loaders, rule.Overrides); err != nil {
			return err
		}
	}

	return nil
}

type NewFSOptions struct {
	Rules []*LoadRule
	// Watch determines if assets should be reloaded on file change.
	Watch bool
}

// NewFS returns a wrapper of an fs.FS with specific loading rules for assets.
// NewImage and NewShader can be used as loaders for *.ext files for example.
func NewFS(fsys fs.FS, opts *NewFSOptions) *FS {
	var watch bool
	loaders := map[string]Loader{}
	if opts != nil {
		watch = opts.Watch
		// embed.FS can't be reloaded
		if _, ok := fsys.(*embed.FS); ok && watch {
			panic("fs: embed.FS is static and cannot be hot-reloaded")
		}
		err := initLoadingRules(fsys, loaders, opts.Rules)
		if err != nil {
			// Malformated file pattern in a rule
			panic(err)
		}
	}
	fmt.Println("loaders:", loaders)

	f := &FS{
		FS: fsys,

		watch:   watch,
		loaders: loaders,
		objects: make(map[string]*fsAsset),
	}
	// Set a finalizer to unregister all files from tracked files at GC
	runtime.SetFinalizer(f, func(f *FS) {
		f.Dispose()
	})

	return f
}

// Get returns the object with the given name, returns an error
// if the file could not be accessed or if asset could not be loaded.
func (fsys *FS) Get(name string) (any, error) {
	if a, ok := fsys.objects[name]; ok {
		return a.obj, nil
	}

	// Read file contents from FS
	data, err := fs.ReadFile(fsys, name)
	if err != nil {
		fmt.Println("there:", name)
		return nil, err
	}
	// Initialize the object
	loader, ok := fsys.loaders[name]
	if !ok {
		return nil, fmt.Errorf("fs: missing loader for '%s'", name)
	}
	obj, err := loader(data)
	if err != nil {
		return nil, err
	}
	a := &fsAsset{
		obj: obj,
		fn:  loader,
	}
	fsys.objects[name] = a
	// Watch the asset if hot reloading is enabled
	if fsys.watch {
		watchAsset(name, a)
	}

	return obj, nil
}

// MustGet returns the object with the given name without error handling.
func (fsys *FS) MustGet(name string) any {
	obj, _ := fsys.Get(name)

	return obj
}

// GetImage works like the Get method but returns the object as an image.
func (fsys *FS) GetImage(name string) (*ebiten.Image, error) {
	obj, err := fsys.Get(name)
	if err != nil {
		return nil, err
	}
	if img, ok := obj.(*ebiten.Image); ok {
		return img, nil
	}

	return nil, fmt.Errorf("fs: '%s' is not an image", name)
}

// MustGetImage returns the image with the given name without error handling.
func (fsys *FS) MustGetImage(name string) *ebiten.Image {
	img, _ := fsys.GetImage(name)

	return img
}

// GetShader works like the Get method but returns the object as a shader.
func (fsys *FS) GetShader(name string) (*ebiten.Shader, error) {
	obj, err := fsys.Get(name)
	if err != nil {
		return nil, err
	}
	if shader, ok := obj.(*ebiten.Shader); ok {
		return shader, nil
	}

	return nil, fmt.Errorf("fs: '%s' is not a shader", name)
}

// MustGetShader returns the shader with the given name without error handling.
func (fsys *FS) MustGetShader(name string) *ebiten.Shader {
	shader, _ := fsys.GetShader(name)

	return shader
}

func (fsys *FS) Dispose() {
	if fsys == nil {
		return
	}
	// Unwatch any tracked file
	for path, a := range fsys.objects {
		unwatchAsset(path, a)
	}
	fsys.loaders = nil
	fsys.objects = nil
	fsys.FS = nil
}
