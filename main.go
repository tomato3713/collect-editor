package main

import (
	"log"
	"os"
	"strings"

	termbox "github.com/nsf/termbox-go"
)

var (
	mode Mode

	editBufs   *Buffer
	cmdLineBuf *Buffer

	cmdLineWin *CmdLineWin
	editWins   *EditWin
)

func main() {
	if err := startUp(); err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

	if len(os.Args) > 1 {
		editBufs.filename = os.Args[1]
	}

	if len(editBufs.filename) <= 0 {
		editBufs.filename = "newfile.txt"
	} else {
		file, err := os.Open(editBufs.filename)
		if err != nil {
			log.Fatal(err)
		}
		editBufs.readFileToBuf(file)
	}

	screenPaint()

	// poll for keyboard events in another goroutine
	events := make(chan termbox.Event, 1000)
	go func() {
		for {
			events <- termbox.PollEvent()
		}
	}()

mainloop:
	for {
		select {
		case ev := <-events:
			switch mode {
			case Move:
				if ev.Type == termbox.EventKey {
					switch ev.Key {
					case termbox.KeyEsc:
						mode = Move
					case termbox.KeyArrowUp:
						editBufs.moveCursor(Up)
					case termbox.KeyArrowDown:
						editBufs.moveCursor(Down)
					case termbox.KeyArrowLeft:
						editBufs.moveCursor(Left)
					case termbox.KeyArrowRight:
						editBufs.moveCursor(Right)
					case termbox.KeyCtrlS:
						editBufs.writeBufToFile()
					case termbox.KeyCtrlC:
						break mainloop // 実行終了
					default:
					}
					switch ev.Ch {
					case ':':
						mode = Cmd
						cmdLineWin.focus()
					case 'k':
						editBufs.moveCursor(Up)
					case 'j':
						editBufs.moveCursor(Down)
					case 'h':
						editBufs.moveCursor(Left)
					case 'l':
						editBufs.moveCursor(Right)
					case 'i':
						mode = Edit
					case 'u':
						editBufs.undo()
					case 'r':
						editBufs.redo()
					default:
					}
				}
			case Edit:
				if ev.Type == termbox.EventKey {
					switch ev.Key {
					case termbox.KeyEsc:
						mode = Move
					case termbox.KeyEnter:
						editBufs.lineFeed()
						// mac delete-key is this
					case termbox.KeyCtrlH:
						fallthrough
					case termbox.KeyBackspace2:
						editBufs.backSpace()
					case termbox.KeyCtrlZ:
						editBufs.undo()
					case termbox.KeyCtrlR:
						editBufs.redo()
					default:
						editBufs.insertChr(ev.Ch)
					}
				}
			case Visual:
				if ev.Type == termbox.EventKey {
					switch ev.Key {
					case termbox.KeyEsc:
						mode = Move
					}
				}
			case Cmd:
				if ev.Type == termbox.EventKey {
					switch ev.Key {
					case termbox.KeyEsc:
						mode = Move
					case termbox.KeyEnter:
						// 入力されたコマンドの解析と実行を開始する
						// quit
						usrCmd := cmdLineBuf.lines[0].text[1:]
						if strings.Compare(string(usrCmd), "q") == 0 || strings.Compare(string(usrCmd), "quit") == 0 {
							break mainloop
						}
					default:
						cmdLineBuf.insertChr(ev.Ch)
					}
				}
			default:
			}
			// when entered any key, redraw buffer
			screenPaint()
		default:
			// Nothing
		}
	}
}

func startUp() error {
	// Initialize terminal window
	err := termbox.Init()
	if err != nil {
		return err
	}
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	termbox.SetCursor(0, 0)

	editBufs = newbuffer()
	cmdLineBuf = newbuffer()

	// get window size
	w, h := termbox.Size()

	// Set command line window default value
	cmdLineWin = newCmdLineWin(w, h, cmdLineBuf)

	// Set editWins default value
	editWins = newEditWin(0, 0, w, h, editBufs)

	mode = Move
	return nil
}

func screenPaint() {
	// clean all window
	termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	editWins.draw()
	cmdLineWin.draw()

	// 現在のモードに合わせて、カーソルを描く
	if mode.equal(Cmd) {
		cmdLineWin.updateCursor()
	} else {
		editWins.updateCursor()
	}

	// update all window
	termbox.Flush()
}
