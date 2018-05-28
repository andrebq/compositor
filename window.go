package main

import (
	"sync"

	"github.com/veandco/go-sdl2/sdl"
)

type (
	window struct {
		sync.Mutex
		sw *sdl.Window

		widgets widgets
	}
)

func (w *window) addWidget(wi ...*widget) {
	w.Lock()
	w.widgets.add(wi...)
	w.Unlock()
}

func newWindow(windows *windowList, title string, w, h int) (*window, error) {
	sw, err := sdl.CreateWindow(
		title,
		sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		w, h,
		0,
	)

	if err != nil {
		return nil, err
	}

	window := windows.addWindow(&window{
		sw:      sw,
		widgets: make(widgets, 0, 10),
	})

	background, err := newWidget(0, 0, w, h)
	if err != nil {
		return nil, err
	}
	background.clear(randomColor())

	leftCorner, err := newWidget(0, 0, w/2, h/2)
	if err != nil {
		return nil, err
	}

	rightCorner, err := newWidget(w/2, h/2, w/2, h/2)
	if err != nil {
		return nil, err
	}

	window.addWidget(background, leftCorner, rightCorner)
	return window, err
}

func (w *window) background() *widget {
	return w.widgets[0]
}

func (w *window) handleKeyUp(k *sdl.KeyUpEvent) {
	println("window.handleKeyUp")
	if k.Keysym.Sym == sdl.K_SPACE {
		w.background().clear(randomColor())
	}
}

func (w *window) draw() {
	dest, err := w.sw.GetSurface()
	if err != nil {
		return
	}

	var invalidatedRect sdl.Rect
	var rects []sdl.Rect
	for _, widget := range w.widgets {
		if !widget.needRepaint(dest) && !widget.rect.HasIntersection(&invalidatedRect) {
			continue
		}
		widget.draw(dest)
		rects = append(rects, widget.rect)
		invalidatedRect = widget.mergeRect(invalidatedRect)
	}

	if len(rects) == 0 {
		return
	}

	w.sw.UpdateSurfaceRects(rects)
}
