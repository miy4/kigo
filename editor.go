package kigo

import (
	"errors"
)

var Quit = errors.New("Quit")

type Editor struct {
	term *Terminal
}

func NewEditor() *Editor {
	term := NewTerminal()
	return &Editor{term}
}

func (editor *Editor) rows() int {
	return int(editor.term.size.Row)
}

func (editor *Editor) cols() int {
	return int(editor.term.size.Col)
}

func (editor *Editor) drawRows() {
	for y := 0; y < editor.rows(); y++ {
		editor.term.writeString("~")
		editor.term.clearLineRight()
		if y < int(editor.rows()) {
			editor.term.writeString("\r\n")
		}
	}
}

func (editor *Editor) refreshScreen() {
	editor.term.hideCursor()
	editor.term.moveCursorToHome()
	editor.drawRows()
	editor.term.moveCursorToHome()
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
	}

	return nil
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
