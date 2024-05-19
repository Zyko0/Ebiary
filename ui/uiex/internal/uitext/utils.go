package uitext

import "image/color"

func ColorAsFloat32RGB(clr color.Color) float32 {
	if clr == nil {
		return 0
	}
	r, g, b, _ := clr.RGBA()
	return float32((r&255)<<16 + (g&255)<<8 + b&255)
}
