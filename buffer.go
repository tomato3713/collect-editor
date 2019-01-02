package main

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"

	termbox "github.com/nsf/termbox-go"
)

const (
	Up = iota
	Down
	Left
	Right
)

type buffer struct {
	cursor   cursor
	lines    []*line
	filename string
	undoBuf  *bufStack
	redoBuf  *bufStack
}

type cursor struct {
	x int
	y int
}

type line struct {
	text []rune
}

type bufStack struct {
	bufs []*buffer
}

func (b *buffer) lineFeed() {
	p := b.cursor.y + 1
	// split line by the cursor and store these
	fh, lh := b.lines[b.cursor.y].split(b.cursor.x)

	t := make([]*line, len(b.lines), cap(b.lines)+1)
	copy(t, b.lines)
	b.lines = append(t[:p+1], t[p:]...)
	b.lines[p] = new(line)

	// write back previous line and newline
	b.lines[p-1].text = fh
	b.lines[p].text = lh

	b.cursor.x = 0
	b.cursor.y++
}

func (b *buffer) backSpace() {
	if b.cursor.x == 0 && b.cursor.y == 0 {
		// nothing to do
	} else {
		if b.cursor.x == 0 {
			// stre current line
			t := b.lines[b.cursor.y].text
			// delete current line
			b.lines = append(b.lines[:b.cursor.y], b.lines[b.cursor.y+1:]...)
			b.cursor.y--
			// join stored lines to previous line-end
			plen := b.lines[b.cursor.y].text
			b.lines[b.cursor.y].text = append(b.lines[b.cursor.y].text, t...)
			b.cursor.x = len(plen)
		} else {
			b.lines[b.cursor.y].deleteChr(b.cursor.x)
			b.cursor.x--
		}
	}
}

func (b *buffer) insertChr(r rune) {
	b.lines[b.cursor.y].insertChr(r, b.cursor.x)
	b.cursor.x++
}

func (l *line) insertChr(r rune, p int) {
	t := make([]rune, len(l.text), cap(l.text)+1)
	copy(t, l.text)
	l.text = append(t[:p+1], t[p:]...)
	l.text[p] = r
}

func (l *line) deleteChr(p int) {
	p = p - 1
	l.text = append(l.text[:p], l.text[p+1:]...)
}

func (b *buffer) updateLines() {
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	for y, l := range b.lines {
		for x, r := range l.text {
			termbox.SetCell(x, y, r, termbox.ColorWhite, termbox.ColorBlack)
		}
	}
}

func (b *buffer) moveCursor(d int) {
	switch d {
	case Up:
		// guard of top "rows"
		if b.cursor.y > 0 {
			b.cursor.y--
			// guard of end of "row"
			if b.cursor.x > len(b.lines[b.cursor.y].text) {
				b.cursor.x = len(b.lines[b.cursor.y].text)
			}
		}
		break
	case Down:
		// guard of end of "rows"
		if b.cursor.y < b.linenum()-1 {
			b.cursor.y++
			// guard of end of "row"
			if b.cursor.x > len(b.lines[b.cursor.y].text) {
				b.cursor.x = len(b.lines[b.cursor.y].text)
			}
		}
		break
	case Left:
		if b.cursor.x > 0 {
			b.cursor.x--
		} else {
			// guard of top of "rows"
			if b.cursor.y > 0 {
				b.cursor.y--
				b.cursor.x = len(b.lines[b.cursor.y].text)
			}
		}
		break
	case Right:
		if b.cursor.x < b.lines[b.cursor.y].runenum() {
			b.cursor.x++
		} else {
			// guard of end of "rows"
			if b.cursor.y < b.linenum()-1 {
				b.cursor.x = 0
				b.cursor.y++
			}
		}
		break
	default:
	}
}

func (b *buffer) updateCursor() {
	termbox.SetCursor(b.cursor.x, b.cursor.y)
}

func (b *buffer) linenum() int {
	return len(b.lines)
}

func (l *line) runenum() int {
	return len(l.text)
}

func (l *line) split(pos int) ([]rune, []rune) {
	return l.text[:pos], l.text[pos:]
}

func (l *line) joint() *line {
	return nil
}

func (b *buffer) pushBufToUndoRedoBuffer() {
	tb := new(buffer)
	tb.cursor.x = b.cursor.x
	tb.cursor.y = b.cursor.y
	for i, l := range b.lines {
		tl := new(line)
		tb.lines = append(tb.lines, tl)
		tb.lines[i].text = l.text
	}
	b.undoBuf.bufs = append(b.undoBuf.bufs, tb)
}

func (b *buffer) undo() {
	if len(b.undoBuf.bufs) == 0 {
		return
	}
	if len(b.undoBuf.bufs) > 1 {
		b.redoBuf.bufs = append(b.redoBuf.bufs, b.undoBuf.bufs[len(b.undoBuf.bufs)-1])
		b.undoBuf.bufs = b.undoBuf.bufs[:len(b.undoBuf.bufs)-1]
	}
	tb := b.undoBuf.bufs[len(b.undoBuf.bufs)-1]
	b.undoBuf.bufs = b.undoBuf.bufs[:len(b.undoBuf.bufs)-1]
	b.cursor.x = tb.cursor.x
	b.cursor.y = tb.cursor.y
	for i, l := range tb.lines {
		tl := new(line)
		b.lines = append(b.lines, tl)
		b.lines[i].text = l.text
	}
}

func (b *buffer) redo() {
	if len(b.redoBuf.bufs) == 0 {
		return
	}
	tb := b.redoBuf.bufs[len(b.redoBuf.bufs)-1]
	b.redoBuf.bufs = b.redoBuf.bufs[:len(b.redoBuf.bufs)-1]
	b.cursor.x = tb.cursor.x
	b.cursor.y = tb.cursor.y
	for i, l := range tb.lines {
		tl := new(line)
		b.lines = append(b.lines, tl)
		b.lines[i].text = l.text
	}
}

func (b *buffer) readFileToBuf(reader io.Reader) error {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		l := new(line)
		l.text = []rune(scanner.Text())
		b.lines = append(b.lines, l)
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

func (b *buffer) writeBufToFile() {
	content := make([]byte, 1024)
	for _, l := range b.lines {
		l.text = append(l.text, '\n')
		content = append(content, string(l.text)...)
	}
	ioutil.WriteFile(b.filename, content, os.ModePerm)
}

func (b *buffer) getline(n int) ([]rune, error) {
	if n < 0 {
		return []rune{}, errors.New("outsideLineNumber")
	}
	if n > b.getlastlinenum() {
		return []rune{}, errors.New("outsideLineNumber")
	}
	return b.lines[n].text, nil
}

func (b *buffer) getlastlinenum() int {
	return len(b.lines)
}

func (b *buffer) redraw() {
	b.updateLines()
	b.updateCursor()
	b.pushBufToUndoRedoBuffer()
	termbox.Flush()
}

func newbuffer() *buffer {
	buf := new(buffer)
	buf.redoBuf = &bufStack{}
	buf.undoBuf = &bufStack{}
	buf.filename = ""
	return buf
}
