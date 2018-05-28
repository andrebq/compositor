package main

import (
	"context"
	"runtime"

	"github.com/sirupsen/logrus"
	"github.com/veandco/go-sdl2/sdl"
)

type (
	action func(ctx context.Context)
)

var (
	contextDoneSDLEventType     uint32
	processActionQueueEventType uint32
	eventQueue                  chan action
	windows                     *windowList
)

func init() {
	runtime.LockOSThread()

	sdl.Init(sdl.INIT_EVERYTHING)

	contextDoneSDLEventType = sdl.RegisterEvents(2)
	processActionQueueEventType = contextDoneSDLEventType + 1
}

func main() {
	eventQueue = make(chan action, 1000)
	windows = newWindowList()

	w, err := newWindow(windows, "hello world", 800, 600)
	if err != nil {
		logrus.WithError(err).Fatal("unable to create window")
	}

	println("window ", w.sw.GetID())

	ctx, cancel := context.WithCancel(withSDLCancel(context.Background()))
	defer cancel()

	loopForever(ctx)
}

func withSDLCancel(ctx context.Context) context.Context {
	go func(ctx context.Context) {
		<-ctx.Done()
		sdl.PushEvent(sdl.UserEvent{
			Type: contextDoneSDLEventType,
		})
	}(ctx)

	return ctx
}

func loopForever(ctx context.Context) {
	for {
		for _, w := range windows.items {
			w.draw()
		}
		for ev := sdl.WaitEvent(); ev != nil; ev = sdl.PollEvent() {
			if _, ok := ev.(*sdl.QuitEvent); ok {
				return
			}
			handleEvent(ctx, ev)
		}
	}
}

func handleEvent(ctx context.Context, ev sdl.Event) {
	switch ev := ev.(type) {
	case *sdl.UserEvent:
		handleUserEvent(ctx, ev)
	case *sdl.KeyUpEvent:
		println("key up")
		win := windows.windowById(ev.WindowID)
		if win != nil {
			win.handleKeyUp(ev)
		}
	}
}

func handleUserEvent(ctx context.Context, ev *sdl.UserEvent) {
	switch ev.Type {
	case contextDoneSDLEventType:
		sdl.Quit()
	case processActionQueueEventType:
		for {
			select {
			case <-ctx.Done():
				return
			case action := <-eventQueue:
				action(ctx)
			default:
				// no more items in the buffer queue
				return
			}
		}
	}
}
