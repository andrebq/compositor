package main

import (
	"image"
	"image/color"

	"math/rand"

	"golang.org/x/image/colornames"
)

func randomColor() color.RGBA {
	name := colornames.Names[int(rand.Int31n(int32(len(colornames.Names))))]
	return colornames.Map[name]
}

func randomImage(w, h int) *image.RGBA {
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			rgba.Set(i, j, randomColor())
		}
	}
	return rgba
}
