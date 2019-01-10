package main

import (
	termbox "github.com/nsf/termbox-go"
)

// BufferWin define window for buffer
type BufferWin struct {
	window
	stsLineHeight int
}

func (w BufferWin) draw() {
	w.updateBufBody()
	w.drawStatusLine()
	w.buf.pushBufToUndoRedoBuffer()
}

func (w BufferWin) drawStatusLine() {
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

func (w BufferWin) focus() {
	w.updateCursor()
}
