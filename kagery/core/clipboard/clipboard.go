//go:build !js

package clipboard

import (
	"golang.design/x/clipboard"
)

func ReadText() string {
	b := clipboard.Read(clipboard.FmtText)
	if len(b) > 0 {
		return string(b)
	}
	return ""
}

func ReadImage() {
	// TODO:
}

func WriteText(txt string) {
	clipboard.Write(clipboard.FmtText, []byte(txt))
}
