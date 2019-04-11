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
	// Cmdモードに入ったときに、cmdline Windowの中身を削除する
	w.buf.deleteLine(1)
	w.buf.insertChr('>')
}

func newCmdLineWin(width int, height int, buf *Buffer) *CmdLineWin {
	w := new(CmdLineWin)
	w.buf = buf
	w.coord.x = 0
	w.coord.y = height - cmdLineHeight
	w.size.width = width
	w.size.height = cmdLineHeight
	return w
}
