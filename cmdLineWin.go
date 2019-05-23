package main

import (
	"github.com/homedm/collect-editor/pkg/buffer"
	"github.com/homedm/collect-editor/pkg/drawing"
	termbox "github.com/nsf/termbox-go"
)

// CmdLineWin define window for command line
type CmdLineWin struct {
	drawing.Window
}

func (w CmdLineWin) Redraw() {
	w.Draw()
	w.DrawStatusLine()
}

func (w CmdLineWin) Focus() {
	// Cmdモードに入ったときに、cmdline Windowの中身を削除する
	w.Buf.DeleteLine(1)
	w.Buf.InsertChr('>')
}

func (w CmdLineWin) DrawStatusLine() {
	// Write filename on status line
	fg := termbox.ColorBlack
	bg := termbox.ColorWhite

	sx := w.Coord.X
	ex := w.Coord.X + w.Size.Width

	ch := []rune("Editing File: " + w.Buf.GetFileName())

	i := 0
	for x := sx; x < ex; x++ {
		if i < len(ch) {
			drawing.DrawChr(x, w.Coord.Y-1, ch[i], fg, bg)
			i++
		} else {
			drawing.DrawChr(x, w.Coord.Y-1, rune(0), fg, bg)
		}
	}
}

// NewCmdLineWin is make new command line window
func NewCmdLineWin(width int, height int, buf *buffer.Buffer) *CmdLineWin {
	w := new(CmdLineWin)
	w.Buf = buf
	_, h := termbox.Size()

	w.Coord.X = 0
	w.Coord.Y = h - height + 1

	w.Size.Width = width
	w.Size.Height = height - 1
	w.Top = 0
	return w
}
