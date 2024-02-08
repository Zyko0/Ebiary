# Parade

Parade (PARAllax+DEpth) is an Ebitengine library to add depth to environments as well as sprites (sprites stacking).

## Usage

`go get github.com/Zyko0/Ebiary/parade`

```go
import "github.com/Zyko0/Ebiary/parade"
```

```go
shader, err := asset.NewLiveAsset[*ebiten.Shader]("myfiles/shader.kage")
if err != nil {
    log.Fatal(err)
}

// ...

if err := shader.Error(); err != nil {
    fmt.Println("err:", err)
}
screen.DrawTrianglesShader(vertices, indices, shader.Value(), nil)
```

If the asset failed to reload, then the method `.Error()` will return a loading error, that, is either due to a failure to read the file on disk, or a failure to parse its content.

When an error shows up, you need to log it explicitely because the content (accessed by `.Value()`) will remain unchanged until the next reload is successful.
This is necessary for shader files, because ebitengine returning a nil shader would crash the program, if your code editor performs auto-save and that you didn't finish updating its code, for example.

## Features

- `*ebiten.Image` => Supports png, jpg, bmp decoding
- `*ebiten.Shader` => It just calls `ebiten.NewShader`

## Notes

- This uses `github.com/fsnotify/fsnotify` to watch for file changes and update the related assets automatically.
- I used generics so that the type assertion does not need to be made on the user side and I also find that it documents the code more, on top of the variable name.
- `Dispose()` does not have to be called.
- Thinking about supporting more arbitrary types from the `audio` package (note that `[]byte` is supported) and custom ones (with loading functions provided by the user)
- Only `*ebiten.Shader` has been tested yet!
- PRs are welcome!