package main

import (
	"log"
	"os"
	"strings"

	"github.com/homedm/collect-editor/pkg/buffer"
	termbox "github.com/nsf/termbox-go"
)

var (
	mode Mode

	editBufs   *buffer.Buffer
	cmdLineBuf *buffer.Buffer

	cmdLineWin *CmdLineWin
	editWins   *EditWin
)

func main() {
	if err := startUp(); err != nil {
		log.Fatal(err)
	}
	defer termbox.Close()

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
						editBufs.MoveCursor(buffer.Up)
					case termbox.KeyArrowDown:
						editBufs.MoveCursor(buffer.Down)
					case termbox.KeyArrowLeft:
						editBufs.MoveCursor(buffer.Left)
					case termbox.KeyArrowRight:
						editBufs.MoveCursor(buffer.Right)
					case termbox.KeyCtrlS:
						editBufs.WriteBufToFile()
					case termbox.KeyCtrlC:
						break mainloop // 実行終了
					default:
					}
					switch ev.Ch {
					case ':':
						mode = Cmd
						cmdLineWin.focus()
					case 'k':
						editBufs.MoveCursor(buffer.Up)
					case 'j':
						editBufs.MoveCursor(buffer.Down)
					case 'h':
						editBufs.MoveCursor(buffer.Left)
					case 'l':
						editBufs.MoveCursor(buffer.Right)
					case 'i':
						mode = Edit
					case 'u':
						editBufs.Undo()
					case 'r':
						editBufs.Redo()
					default:
					}
				}
			case Edit:
				if ev.Type == termbox.EventKey {
					switch ev.Key {
					case termbox.KeyEsc:
						mode = Move
					case termbox.KeyEnter:
						editBufs.LineFeed()
						// mac delete-key is this
					case termbox.KeyCtrlH:
						fallthrough
					case termbox.KeyBackspace2:
						editBufs.BackSpace()
					case termbox.KeyCtrlZ:
						editBufs.Undo()
					case termbox.KeyCtrlR:
						editBufs.Redo()
					default:
						editBufs.InsertChr(ev.Ch)
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
						usrCmd, err := cmdLineBuf.GetLine(0)
						if err != nil {
							log.Fatal(err)
						}
						if strings.Compare(string(usrCmd[1:]), "q") == 0 || strings.Compare(string(usrCmd), "quit") == 0 {
							break mainloop
						}
					default:
						cmdLineBuf.InsertChr(ev.Ch)
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

	if len(os.Args) > 1 {
		fname := os.Args[1]
		editBufs, err = buffer.NewBuffer(fname)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		editBufs = buffer.NewEmptyBuffer()
	}

	editBufs = buffer.NewEmptyBuffer()
	cmdLineBuf = buffer.NewEmptyBuffer()

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
		cmdLineWin.UpdateCursor()
	} else {
		editWins.UpdateCursor()
	}

	// update all window
	termbox.Flush()
}
