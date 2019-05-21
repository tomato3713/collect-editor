package main

import (
	"github.com/homedm/collect-editor/pkg/buffer"
	"github.com/homedm/collect-editor/pkg/drawing"
	termbox "github.com/nsf/termbox-go"
)

// EditWin define window for buffer
type EditWin struct {
	drawing.Window
	stsLineHeight int
}

func (w EditWin) draw() {
	w.UpdateBufBody()
	w.drawStatusLine()
	w.Buf.PushBufToUndoRedoBuffer()
}

func (w EditWin) drawStatusLine() {
	// Write filename on status line
	fg := termbox.ColorBlack
	bg := termbox.ColorWhite

	sx := w.Coord.X
	ex := w.Coord.X + w.Size.Width

	ch := []rune(w.Buf.GetFileName())

	i := 0
	for x := sx; x < ex; x++ {
		if i < len(ch) {
			termbox.SetCell(x, w.Coord.Y+w.Size.Height-w.stsLineHeight, ch[i], fg, bg)
			i++
		} else {
			termbox.SetCell(x, w.Coord.Y+w.Size.Height-w.stsLineHeight, rune(0), fg, bg)
		}
	}
}

func (w EditWin) focus() {
	w.UpdateCursor()
}

func newEditWin(x int, y int, width int, height int, buf *buffer.Buffer) *EditWin {
	w := new(EditWin)
	w.Buf = buf
	w.Coord.X = x
	w.Coord.Y = y
	w.Size.Width = width
	w.Size.Height = height - cmdLineHeight
	w.stsLineHeight = 1
	return w
}
