package main

import (
	termbox "github.com/nsf/termbox-go"
)

// EditWin define window for buffer
type EditWin struct {
	window
	stsLineHeight int
}

func (w EditWin) draw() {
	w.updateBufBody()
	w.drawStatusLine()
	w.buf.pushBufToUndoRedoBuffer()
}

func (w EditWin) drawStatusLine() {
	// Write filename on status line
	fg := termbox.ColorBlack
	bg := termbox.ColorWhite

	sx := w.coord.x
	ex := w.coord.x + w.size.width

	ch := []rune(w.buf.filename)

	i := 0
	for x := sx; x < ex; x++ {
		if i < len(ch) {
			termbox.SetCell(x, w.coord.y+w.size.height-w.stsLineHeight, ch[i], fg, bg)
			i++
		} else {
			termbox.SetCell(x, w.coord.y+w.size.height-w.stsLineHeight, rune(0), fg, bg)
		}
	}
}

func (w EditWin) focus() {
	w.updateCursor()
}

func newEditWin(width int, height int, buf *Buffer) *EditWin {
	w := new(EditWin)
	w.buf = buf
	w.coord.x = 0
	w.coord.y = 0
	w.size.width = width
	w.size.height = height - cmdLineHeight
	w.stsLineHeight = 1
	return w
}
