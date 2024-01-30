# JFA

An implementation of the [jump flooding algorithm (JFA)](https://en.wikipedia.org/wiki/Jump_flooding_algorithm) using Ebitengine.
It can be used to create signed distance field textures, with fonts for example.

Note that there are still a few edge cases that need to be fixed overall!

## Usage

`go get github.com/Zyko0/Ebiary/jfa`

```go
import "github.com/Zyko0/Ebiary/jfa"
```

```go
// Create a destination image to store the computed distance info 
dstImg := ebiten.NewImage(width, height)
// Create a jfa instance with the desired image size
j := jfa.New(dstImg.Bounds())
// Compute the distance of the given "maskImg" ebiten image
j.Generate(dstImg, maskImg, &jfa.GenerateOptions{
    // Pixel colors with total luminance (grey) > 0 and alpha > 0 will
    // contribute to the plain shape
    PlainValueThresholds: map[jfa.ColorMask]float64{
        jfa.ColorMaskAlpha:     0,
        jfa.ColorMaskGreyscale: 0,
    },
    // The same distance value will be set in R,G,B channels
    Encoding:     jfa.EncodingDistanceGreyscale,
    // The exterior distance to the shape is calculated
    DistanceType: jfa.DistanceExterior,
    // The edges do not count as part of the plain shape
    EdgesPlain:   false,
    Steps:        1024,
    JumpDistance: 8,
})
```

## Previews

| Original | Exterior |
|-|-|
|![img](https://github.com/Zyko0/Ebiary/assets/13394516/09a2257f-ebd5-428f-8691-e80e227167e1)|![image](https://github.com/Zyko0/Ebiary/assets/13394516/76e136c7-2d62-424c-b747-023b7d738562)|

| Exterior (Edges plain) | Interior |
|-|-|
|![image](https://github.com/Zyko0/Ebiary/assets/13394516/f6e5d3ef-1c50-435c-90cc-e46ec9613f46)|![image](https://github.com/Zyko0/Ebiary/assets/13394516/d0366527-9d95-4479-856d-b9b15fe351fc)|

| UVs (Exterior) | UVs (Interior) |
|-|-|
|![image](https://github.com/Zyko0/Ebiary/assets/13394516/781c918a-4729-42ac-88b3-3f2fc90d9abe)|![image](https://github.com/Zyko0/Ebiary/assets/13394516/e528798f-3a1a-4400-997d-700e965674ce)|


## CLI (TBD)

Not done yet

`go run ./cmd "test.png"` / `go run ./cmd` (drag & drop)

## Notes

- This is an early version, a few things are not clearly defined yet
  - Different kinds of encoding methods for example
  - There's a bug with the default values around the edges that I don't know how to fix
- `Steps` and `JumpDistance` must be played with in order to get satisfying results, for more information about the theory I recommend to look online!
- Any PR contribution (fix/improvement/feature) is welcome btw!
