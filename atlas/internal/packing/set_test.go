package packing

import (
	"image"
	"math/rand"
	"testing"
)

func TestSet_Insert(t *testing.T) {
	s := NewSet(2048, 2048, nil)
	t.Run("Fill some", func(t *testing.T) {
		for i := 0; i < 4096; i++ {
			r := image.Rect(0, 0, 8+rand.Intn(8), 8+rand.Intn(8))
			s.Insert(&r)
		}
	})
}
