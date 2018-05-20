package main

import (
	"context"
	"image/color"
	"sync"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type (
	window struct {
		sync.Mutex
		pixw   *pixelgl.Window
		ctx    context.Context
		cancel context.CancelFunc

		bg *Image
	}
)

func createWindow(ctx context.Context, title string, width, height int) (*window, error) {
	cfg := pixelgl.WindowConfig{
		Title:  title,
		Bounds: pixel.R(0, 0, float64(width), float64(height)),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	bg := RandomImage(width, height)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(ctx)
	return &window{
		pixw:   win,
		bg:     bg,
		ctx:    ctx,
		cancel: cancel,
	}, nil
}

func (w *window) Clear(rgba color.RGBA) {
	w.pixw.Clear(rgba)
}

func (w *window) Closed() bool {
	return w.pixw.Closed()
}

func (w *window) Update() {
	w.Lock()
	defer w.Unlock()

	w.bg.pic.Draw(w.pixw, pixel.IM.Moved(w.pixw.Bounds().Center()))
	w.pixw.Update()
}

func (w *window) ChangeBG() {
	w.Lock()
	defer w.Unlock()

	bg := RandomImage(int(w.pixw.Bounds().W()),
		int(w.pixw.Bounds().H()))
	w.bg = bg
}

func run(ctx context.Context) {

	win, err := createWindow(ctx, "pixel rocks", 1024, 768)
	if err != nil {
		panic(err)
	}

	win.Clear(colornames.Skyblue)

	go func() {
		for {
			ticker := time.NewTicker(time.Millisecond * 16)
			defer ticker.Stop()
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				win.ChangeBG()
			}
		}
	}()

	for !win.Closed() {
		win.Update()
	}
}

func main() {
	ctx := context.Background()
	pixelgl.Run(func() {
		run(ctx)
	})
}
