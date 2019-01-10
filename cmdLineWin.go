package main

const (
	cmdLineHeight = 2
)

// CmdLineWin define window for command line
type CmdLineWin struct {
	window
}

func (w CmdLineWin) draw() {
	w.updateBufBody()
	w.buf.pushBufToUndoRedoBuffer()
}

func (w CmdLineWin) focus() {
	w.buf.deleteLine(1)
	w.updateBufBody()
	w.updateCursor()
	w.buf.insertChr('>')
}
