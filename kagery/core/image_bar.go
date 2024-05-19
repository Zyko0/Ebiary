package core

import (
	"image"
	"io/fs"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	"github.com/Zyko0/Ebiary/kagery/assets"
	"github.com/Zyko0/Ebiary/ui"
	"github.com/Zyko0/Ebiary/ui/opt"
	"github.com/Zyko0/Ebiary/ui/uiex"
	"github.com/hajimehoshi/ebiten/v2"
)

type ImageBar struct {
	*ui.Grid

	pictures [4]*uiex.Picture
}

func newPicture(logger *Logger, img *ebiten.Image) *uiex.Picture {
	return uiex.NewPicture(img).WithOptions(
		opt.Picture.Image.Options(
			opt.Image.FillContainer(true),
		),
		opt.Picture.Options(
			ui.WithCustomUpdateFunc(func(pic *uiex.Picture, is ui.InputState) {
				if !pic.Hovered() {
					return
				}

				files := is.DroppedFiles()
				if files == nil {
					return
				}
				var img *ebiten.Image
				var imgPath string
				err := fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
					if d.IsDir() {
						return nil
					}
					if err != nil {
						return err
					}
					if img != nil {
						return nil
					}
					imgPath = path
					f, err := files.Open(path)
					if err != nil {
						return err
					}
					raw, _, err := image.Decode(f)
					if err != nil {
						return err
					}
					img = ebiten.NewImageFromImage(raw)
					return nil
				})
				if err != nil {
					logger.RegisterError(&imageError{
						errorBase: &errorBase{
							msg: err.Error() + ":" + imgPath,
						},
					})
					return
				}
				if img != nil {
					pic.Image().Image().Deallocate()
					pic.Image().SetImage(img)
				}
			}),
		),
	)
}

func NewImageBar(logger *Logger) *ImageBar {
	ib := &ImageBar{
		Grid: ui.NewGrid(1, 4).WithOptions(
			opt.Grid.Options(
				opt.RGB(0, 0, 0),
				opt.Alpha(Alpha),
				opt.Rounding(Rounding),
				opt.Padding(8),
				opt.Margin(-8),
				// Hack: Extra -2 margin for the border not to be drawn
				// over by scrollbars
				opt.MarginBottom(-10),
				opt.MarginTop(-10),
				opt.Border(1, clrBorder),
			),
		),
		pictures: [4]*uiex.Picture{
			newPicture(logger, assets.GopherImage),
			newPicture(logger, assets.NormalImage),
			newPicture(logger, assets.GopherBgImage),
			newPicture(logger, assets.NoiseImage),
		},
	}
	ib.Add(0, 0, 1, 1, ib.pictures[0])
	ib.Add(0, 1, 1, 1, ib.pictures[1])
	ib.Add(0, 2, 1, 1, ib.pictures[2])
	ib.Add(0, 3, 1, 1, ib.pictures[3])

	return ib
}

func (ib *ImageBar) Images() [4]*ebiten.Image {
	return [4]*ebiten.Image{
		ib.pictures[0].Image().Image(),
		ib.pictures[1].Image().Image(),
		ib.pictures[2].Image().Image(),
		ib.pictures[3].Image().Image(),
	}
}
