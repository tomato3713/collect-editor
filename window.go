package main

import (
	termbox "github.com/nsf/termbox-go"
)

type window struct {
	// The coordinates of the upper left corner
	coord coordinal
	size  size
	buf   *Buffer
}

type coordinal struct {
	x int
	y int
}

type size struct {
	width  int
	height int
}

func (w window) drawHorizon(sx int, sy int, ex int, ey int, fg termbox.Attribute, bg termbox.Attribute) {
	var ch rune
	for x := sx; x < ex; x++ {
		for y := sy; y < ey; y++ {
			termbox.SetCell(x, y, ch, fg, bg)
		}
	}
}

func (w window) updateCursor() {
	termbox.SetCursor(w.buf.cursor.x+w.coord.x, w.buf.cursor.y+w.coord.y)
}

func (w window) updateBufBody() {
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
