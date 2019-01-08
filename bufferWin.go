package main

import (
	termbox "github.com/nsf/termbox-go"
)

type BufferWin struct {
	window
	buf *Buffer
}

func (w BufferWin) draw() {
	w.updateBufBody()
	w.updateCursor()
	w.drawStatusLine()
	w.buf.pushBufToUndoRedoBuffer()
}

func (w BufferWin) updateCursor() {
	termbox.SetCursor(w.buf.cursor.x+w.coord.x, w.buf.cursor.y+w.coord.y)
}

func (w BufferWin) updateBufBody() {
	// Clear inside this window
	for y := w.coord.y; y < w.coord.y+w.size.height; y++ {
		for x := w.coord.x; x < w.coord.x+w.size.width; x++ {
			termbox.SetCell(y, x, rune(0), termbox.ColorWhite, termbox.ColorBlack)
		}
	}
	// Draw text inside this window
	for y, l := range w.buf.lines {
		for x, r := range l.text {
			termbox.SetCell(x+w.coord.x, y+w.coord.y, r, termbox.ColorWhite, termbox.ColorBlack)
		}
	}
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
			termbox.SetCell(x, w.coord.y+w.size.height, ch[i], fg, bg)
			i++
		} else {
			termbox.SetCell(x, w.coord.y+w.size.height, rune(0), fg, bg)
		}
	}
}
