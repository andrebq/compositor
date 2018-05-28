package main

type (
	widgets []*widget
)

func (ws *widgets) add(w ...*widget) {
	*ws = append(*ws, w...)
	println("ws", *ws, "w", w)
}
