# Atlas

A small package providing an `Atlas` type as well as `atlas.Image` which can be used for more control over batching draw calls.

It works by alleviating some generic work done by ebitengine when treating `DrawImage` commands.
Here, since the images are just regions of an atlas image, the commands can be converted to a slice of triangles that can be submitted in a single ebitengine command, relieving ebitengine from doing additional unnecessary work.

Originally made to address this: https://github.com/hajimehoshi/ebiten/issues/2976

## Usage

`go get github.com/Zyko0/Ebiary/atlas`

- Empty atlas 
```go
// Create an atlas
atlas := atlas.New(1024, 1024, nil)
// Allocate an image on the atlas at a random location
img := app.atlas.NewImage(48, 48)
// Access the *ebiten.Image as a sub-image and set its content
img.Image().Fill(color.White)
```
- Atlas from spritesheet
```go
// Create an atlas
atlas := atlas.New(1024, 1024, nil)
// Set the atlas image's content with the spritesheet image content
atlas.Image().DrawImage(spriteSheet, nil)
// Locate the sprite sub images on the atlas manually
img := app.atlas.SubImage(image.Rect(0,0,48,48))
```
- DrawList
```go
// It's best to re-use the DrawList as Flush() clears the triangles
// but keeps the capacities for next Add commands
dl := &atlas.DrawList{}
dc := &atlas.DrawCommand{}
for _, s := range sprites {
	dc.Image = s.Image
	dc.GeoM = s.GeoM
	dl.Add(dc)
}
// Flush all draws to the screen
dl.Flush(screen, nil)
```

An example can be found in [examples/usage/main.go](./examples/usage/main.go)

### Notes

- The `Text*` related stuff is not implemented properly and will probably never work, since I realized we lack context from `text/v2` to cache all the variations of a glyphs. This results in a drawn text but some characters will have an incorrect X offset. 
- The `Free` method is not implemented yet either at the moment, the `Atlas` is meant to be written to, and refered to by `atlas.Image` but not meant as a dynamic object capable of deallocations yet.
