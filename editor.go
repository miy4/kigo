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

func (editor *Editor) refreshScreen() {
	editor.term.clearEntireScreen()
	editor.term.moveCursorToHome()
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
