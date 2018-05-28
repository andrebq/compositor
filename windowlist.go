package main

import "sync"

type (
	windowList struct {
		sync.Mutex

		items windowMap
	}

	windowMap map[uint32]*window
)

func newWindowList() *windowList {
	return &windowList{
		items: make(windowMap),
	}
}

func (wl *windowList) addWindow(w *window) *window {
	wl.items[w.sw.GetID()] = w
	return w
}

func (wl *windowList) removeWindow(id uint32) *window {
	value := wl.items[id]
	if value != nil {
		return value
	}
	return nil
}

func (wl *windowList) windowById(id uint32) *window {
	for k, w := range wl.items {
		if k == id {
			return w
		}
	}

	return nil
}
