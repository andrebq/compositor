package main

import (
	"image"
	"image/color"

	"math/rand"

	"golang.org/x/image/colornames"

	"github.com/faiface/pixel"
)

type (
	Image struct {
		pic *pixel.Sprite
	}
)

func randomColor() color.RGBA {
	name := colornames.Names[int(rand.Int31n(int32(len(colornames.Names))))]
	return colornames.Map[name]
}

func RandomImage(w, h int) *Image {
	rgba := image.NewRGBA(image.Rect(0, 0, w, h))
	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			rgba.Set(i, j, randomColor())
		}
	}

	pic := pixel.PictureDataFromImage(rgba)
	return &Image{
		pic: pixel.NewSprite(pic, pixel.R(0, 0, float64(w), float64(h))),
	}
}
