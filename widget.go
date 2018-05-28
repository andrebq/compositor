package main

import (
	"image"
	"image/color"
	"sync/atomic"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

type (
	widget struct {
		id   uint64
		surf *sdl.Surface
		rect sdl.Rect

		lastPaint *sdl.Surface
	}
)

var (
	widgetIDs uint64
)

func newWidget(x, y, width, height int) (*widget, error) {
	w := &widget{}
	w.id = atomic.AddUint64(&widgetIDs, 1)

	var err error

	w.surf, err = sdl.CreateRGBSurface(0, int32(width), int32(height), 32, 0, 0, 0, 0)
	if err != nil {
		return nil, err
	}
	w.rect = sdl.Rect{
		X: int32(x),
		Y: int32(y),
		W: w.surf.ClipRect.W,
		H: w.surf.ClipRect.H,
	}

	rndImg := randomImage(width, height)
	w.fillWith(rndImg)

	return w, nil
}

func (w *widget) draw(dest *sdl.Surface) {
	w.surf.Blit(&w.surf.ClipRect, dest, &w.rect)
	w.lastPaint = dest
}

func (w *widget) clear(color color.Color) {
	w.surf.FillRect(&w.surf.ClipRect, toSDLColor(color).Uint32())
	w.dirty()
}

func toSDLColor(c color.Color) sdl.Color {
	r, g, b, a := c.RGBA()
	return sdl.Color{
		R: uint8(r / 255),
		G: uint8(g / 255),
		B: uint8(b / 255),
		A: uint8(a / 255),
	}
}

func (w *widget) fillWith(img *image.RGBA) error {
	surf, err := sdl.CreateRGBSurfaceFrom(unsafe.Pointer(&img.Pix[0]),
		img.Bounds().Dx(),
		img.Bounds().Dy(),
		32, 4*img.Bounds().Dx(), 0, 0, 0, 0)
	if err != nil {
		return err
	}
	defer surf.Free()
	surf.Blit(&surf.ClipRect, w.surf, &w.surf.ClipRect)
	w.dirty()
	return nil
}

func (w *widget) dirty() {
	w.lastPaint = nil
}

func (w *widget) needRepaint(dest *sdl.Surface) bool {
	return w.lastPaint != dest
}

func (w *widget) mergeRect(other sdl.Rect) sdl.Rect {
	return w.rect.Union(&other)
}
