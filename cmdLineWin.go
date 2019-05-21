package main

import (
	"github.com/homedm/collect-editor/pkg/buffer"
	"github.com/homedm/collect-editor/pkg/drawing"
)

const (
	cmdLineHeight = 2
)

// CmdLineWin define window for command line
type CmdLineWin struct {
	drawing.Window
}

func (w CmdLineWin) draw() {
	w.UpdateBufBody()
	w.Buf.PushBufToUndoRedoBuffer()
}

func (w CmdLineWin) focus() {
	// Cmdモードに入ったときに、cmdline Windowの中身を削除する
	w.Buf.DeleteLine(1)
	w.Buf.InsertChr('>')
}

func newCmdLineWin(width int, height int, buf *buffer.Buffer) *CmdLineWin {
	w := new(CmdLineWin)
	w.Buf = buf
	w.Coord.X = 0
	w.Coord.Y = height - cmdLineHeight
	w.Size.Width = width
	w.Size.Height = cmdLineHeight
	return w
}
