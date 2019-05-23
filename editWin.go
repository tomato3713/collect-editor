package main

import (
	"github.com/homedm/collect-editor/pkg/buffer"
	"github.com/homedm/collect-editor/pkg/drawing"
)

// EditWin define window for buffer
type EditWin struct {
	drawing.Window
}

func NewEditWin(x int, y int, width int, height int, buf *buffer.Buffer) *EditWin {
	w := new(EditWin)
	w.Buf = buf
	w.Coord.X = x
	w.Coord.Y = y
	w.Size.Width = width
	w.Size.Height = height - 1
	w.Top = 0
	return w
}
