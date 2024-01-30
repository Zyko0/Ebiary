# JFA

An implementation of the [jump flooding algorithm (JFA)](https://en.wikipedia.org/wiki/Jump_flooding_algorithm) using Ebitengine.
It can be used to create signed distance field textures, with fonts for example.

Note that it's a bit broken at the moment!

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



## CLI (TBD)

Not done yet

`go run ./cmd "test.png"` / `go run ./cmd` (drag & drop)

## Notes

- This is an early version, a few things are not clearly defined yet
  - Different kinds of encoding methods for example
  - There's a bug with the default values around the edges that I don't know how to fix
- `Steps` and `JumpDistance` must be played with in order to get satisfying results, for more information about the theory I recommend to look online!
- Any PR contribution (fix/improvement/feature) is welcome btw!