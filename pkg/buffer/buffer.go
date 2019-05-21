package buffer

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

const (
	// Up is reserved words for commands to move the cursor upward.
	Up = iota
	// Down is reserved words for commands to move the cursor downward.
	Down
	// Left is reserved words for commands to move the cursor to the left.
	Left
	// Right is reserved words for commands to move the cursor to the right.
	Right
)

var (
	// ErrOutRange means variable in out of range
	ErrOutRange = errors.New("error - variable out of range")
	// ErrDeleteProhibitedLine means this line is delete prohibited line
	ErrDeleteProhibitedLine = errors.New("error - this line is delete prohibited line")
)

// Buffer defines buffer structure
type Buffer struct {
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
	bufs []*Buffer
}

func (b *Buffer) LineFeed() {
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

func (b *Buffer) BackSpace() {
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

func (b *Buffer) DeleteLine(n int) error {
	// Check n
	if n <= 0 || n > b.GetLastLineNum() {
		return ErrOutRange
	}

	if b.GetLastLineNum() == 1 {
		b.lines[0].text = []rune{}
		b.cursor.y = 0
		b.cursor.x = 0
	} else {
		b.lines = append(b.lines[:n], b.lines[n+1:]...)
		// update cursor position if cursor is in out of buffer lines
		if b.cursor.y > n {
			b.cursor.y--
		}
		if b.cursor.x > len(b.lines[b.cursor.y].text) {
			b.cursor.x = len(b.lines[b.cursor.y].text)
		}
	}
	return nil
}

func (b *Buffer) InsertChr(r rune) {
	b.lines[b.cursor.y].InsertChr(r, b.cursor.x)
	b.cursor.x++
}

func (l *line) InsertChr(r rune, p int) {
	t := make([]rune, len(l.text), cap(l.text)+1)
	copy(t, l.text)
	l.text = append(t[:p+1], t[p:]...)
	l.text[p] = r
}

func (l *line) deleteChr(p int) {
	p = p - 1
	l.text = append(l.text[:p], l.text[p+1:]...)
}

func (b *Buffer) MoveCursor(d int) {
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

func (b *Buffer) linenum() int {
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

func (b *Buffer) PushBufToUndoRedoBuffer() {
	tb := new(Buffer)
	tb.cursor.x = b.cursor.x
	tb.cursor.y = b.cursor.y
	for i, l := range b.lines {
		tl := new(line)
		tb.lines = append(tb.lines, tl)
		tb.lines[i].text = l.text
	}
	b.undoBuf.bufs = append(b.undoBuf.bufs, tb)
}

func (b *Buffer) Undo() {
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

func (b *Buffer) Redo() {
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

func (b *Buffer) readFileToBuf(reader io.Reader) error {
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

func (b *Buffer) WriteBufToFile() {
	content := make([]byte, 1024)
	for _, l := range b.lines {
		l.text = append(l.text, '\n')
		content = append(content, string(l.text)...)
	}
	ioutil.WriteFile(b.filename, content, os.ModePerm)
}

func (b *Buffer) GetLine(n int) ([]rune, error) {
	if n < 0 {
		return []rune{}, ErrOutRange
	}
	if n > b.GetLastLineNum() {
		return []rune{}, ErrOutRange
	}
	return b.lines[n].text, nil
}

func (b *Buffer) GetLines() []*line {
	return b.lines
}

func (b *Buffer) GetLastLineNum() int {
	return len(b.lines)
}

func (b *Buffer) GetCursor() (x int, y int) {
	x = b.cursor.x
	y = b.cursor.y
	return
}

func (b *Buffer) GetFileName() string {
	return b.filename
}

func NewBuffer(fname string) (*Buffer, error) {
	buf := new(Buffer)
	buf.redoBuf = &bufStack{}
	buf.undoBuf = &bufStack{}
	buf.lines = []*line{&line{[]rune{}}}
	buf.lines[0].text = make([]rune, 0, 0)
	buf.cursor.x = 0
	buf.cursor.y = 0
	buf.filename = fname

	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	buf.readFileToBuf(file)
	return buf, nil
}

func NewEmptyBuffer() *Buffer {
	buf := new(Buffer)
	buf.redoBuf = &bufStack{}
	buf.undoBuf = &bufStack{}
	buf.lines = []*line{&line{[]rune{}}}
	buf.lines[0].text = make([]rune, 0, 0)
	buf.cursor.x = 0
	buf.cursor.y = 0
	return buf
}
