# Asset

A small utility package to ease the management of game assets using Ebitengine, made with development environment in mind mostly.

There should not be any issue (that I'm aware of), but I still recommend to use this library as an helper to start or develop your project / kage shaders (for the live reload feature), and try to get rid of it for your release build.

## Hot reloading assets

`go get github.com/Zyko0/Ebiary/asset`

```go
import "github.com/Zyko0/Ebiary/asset"
```

```go
var (
    img *ebiten.Image
    shader *ebiten.Shader
    // Custom type
    config *Config
)

func init() {
    var err error

    // Loading func is infered automatically for the image type.
    img, err = asset.NewLiveAsset[*ebiten.Image]("myfiles/sprite.png")
    if err != nil {
        log.Fatal(err)
    }
    // Loading func is infered automatically for the shader type.
    shader, err = asset.NewLiveAsset[*ebiten.Shader]("myfiles/shader.kage")
    if err != nil {
        log.Fatal(err)
    }
    // Providing a custom loading function is also possible.
    // For example if you had a Sprite, GameLevel or a Config file to 
    // automatically parse and initialize on file change.
    config, err = NewLiveAssetFunc[*Config]("myfiles/config.json", func (data []byte) (*Config, error) {
        cfg := &Config{}
        err := json.Unmarshal(data, cfg)
        if err != nil {
            return nil, err
        }
        return cfg, nil
    })
}

// ...

func (g *Game) Draw(screen *ebiten.Image) {
    // Log potential reloading errors if necessary (like a kage error).
    // Errors do not nullify the objects accessed by .Value(), instead
    // the previous no-error value is kept, so logging is necessary in case
    // you want to know the impact of your change to the file.
    if img.Error() != nil {
        fmt.Println("warn: image reloading error:", img.Error())
    }
    if shader.Error() != nil {
        fmt.Println("warn: shader reloading error:", shader.Error())
    }

    // Access the most up-to-date value of both 'shader' and 'img' objects.
    screen.DrawRectShader(512, 512, shader.Value(), &ebiten.DrawRectShaderOptions{
        Images: [4]*ebiten.Image{
            img.Value(),
        },
    })
}
```

- The specified type can implement the `Loadable` interface by having its own `func (o *MyObject) Deserialize([]byte) error` function.
- If the asset failed to reload, then the method `.Error()` will return a loading error, that, is either due to a failure to read the file on disk, or a failure to load its content.
- When an error shows up, you need to log it explicitely because the content (accessed by `.Value()`) will remain unchanged until the next reload is successful.
This is necessary for shader files for example, because ebitengine returning a `nil` shader would crash the program.

## FS

FS is a wrapper to `fs.FS` meant to help loading ebitengine assets during development process.

- With a static `embed.FS` (no hot reload)
```go
var (
    //go:embed assets/*.png
    //go:embed assets/shaders/*.kage
    //go:embed assets/config/*.json
    fsys embed.FS

    assetfs = asset.NewFS(fsys, &asset.NewFSOptions{
        Rules: []*asset.LoadRule{
            {
                FilePattern: "asset/*.png",
                Func: asset.NewImage,
                Overrides: []*asset.LoadRule{
                    {
                        FilePattern: "asset/special_sprite.png",
                        Func: asset.NewImageWithOptionsLoader(&ebiten.NewImageOptions{
                            Unmanaged: true,
                        }),
                    },
                },
            },
            {
                FilePattern: "asset/shaders/*.kage",
                Func: asset.NewShader,
            },
            {
                FilePattern: "asset/config/settings.json",
                Func: func(b []byte) (any, error) {
                    var cfg Config
                    err := json.Unmarshal(b, &cfg)
                    if err != nil {
                        return nil, err
                    }
                    return &cfg, nil
                },
            },
        },
    })
)
```
- With `os.DirFS()`, you can replace your above setup with the one below (hot reload possible):
```go
var (
    fsys = os.DirFS(".")

    assetfs = asset.NewFS(fsys, &asset.NewFSOptions{
        Rules: []*asset.LoadRule{
            {
                FilePattern: "asset/*.png",
                Func: asset.NewImage,
            },
            {
                FilePattern: "asset/shaders/*.kage",
                Func: asset.NewShader,
            },
        },
        // Specify that objects should be reloaded on file change
        Watch: true,
    })
)
```
- To use
```go
img := assetfs.MustGetImage("asset/sprite.png")
shader := assetfs.MustGetShader("asset/shaders/blur.kage")
cfg := assetfs.MustGet("asset/config/settings.json").(*Config)
```

### Notes

- This uses `github.com/fsnotify/fsnotify` to watch for file changes and update the related assets automatically.
- I used generics so that the type assertion does not need to be made on the user side and I also find that it documents the code more, on top of the variable name.
- `Dispose()` does not have to be called.