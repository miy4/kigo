package kigo

import (
	"errors"
	"fmt"
	"strings"
)

const version = "0.1.0"

var Quit = errors.New("Quit")

type Editor struct {
	term *Terminal
	cur  *pos
}

type pos struct {
	x int
	y int
}

func NewEditor() *Editor {
	term := NewTerminal()
	return &Editor{
		term: term,
		cur:  &pos{0, 0},
	}
}

func (editor *Editor) rows() int {
	return int(editor.term.size.Row)
}

func (editor *Editor) cols() int {
	return int(editor.term.size.Col)
}

func (editor *Editor) drawRows() {
	for y := 0; y < editor.rows(); y++ {
		if y == editor.rows()/3 {
			editor.drawWelcome()
		} else {
			editor.term.writeString("~")
		}

		editor.term.clearLineRight()

		if y < int(editor.rows()) {
			editor.term.writeString("\r\n")
		}
	}
}

func (editor *Editor) drawWelcome() {
	welcome := fmt.Sprintf("Kigo editor -- version %s", version)
	if len(welcome) > editor.cols() {
		welcome = welcome[:editor.cols()]
	}

	padding := (editor.cols() - len(welcome)) / 2
	if padding > 0 {
		editor.term.writeString("~")
		editor.term.writeString(strings.Repeat(" ", padding-1))
	}

	editor.term.writeString(welcome)
}

func (editor *Editor) refreshScreen() {
	editor.term.hideCursor()
	editor.term.moveCursorToHome()
	editor.drawRows()
	editor.term.moveCursor(editor.cur.y, editor.cur.x)
	editor.term.showCursor()
	editor.term.flush()
}

func (editor *Editor) processKeypress() error {
	key, err := editor.term.ReadKey()
	if err != nil {
		return err
	}

	switch key {
	case KeyCtrlQ:
		return Quit
	case 'w', 's', 'a', 'd':
		editor.moveCursor(key)
	}

	return nil
}

func (editor *Editor) moveCursor(key Key) {
	switch key {
	case 'a':
		editor.cur.x--
	case 'd':
		editor.cur.x++
	case 'w':
		editor.cur.y--
	case 's':
		editor.cur.y++
	}
}

func (editor *Editor) Run() error {
	if err := editor.term.EnableRawMode(); err != nil {
		return err
	}
	defer editor.term.DisableRawMode()

	if err := editor.term.init(); err != nil {
		return err
	}

	for {
		editor.refreshScreen()
		if err := editor.processKeypress(); err != nil {
			editor.term.clearEntireScreen()
			editor.term.moveCursorToHome()

			if err == Quit {
				break
			}

			return err
		}
	}

	return nil
}
