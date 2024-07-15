package packing

import (
	"image"
	"math/rand"
	"testing"
)

func TestSet_Insert(t *testing.T) {
	s := NewSet(2048, 2048, nil)
	t.Run("Fill some", func(t *testing.T) {
		for i := 0; i < 512; i++ {
			r := image.Rect(0, 0, 32+rand.Intn(32), 32+rand.Intn(32))
			s.Insert(&r)
		}
	})
}
