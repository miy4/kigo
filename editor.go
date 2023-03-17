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
		err := editor.processKeypress()
		if err == Quit {
			break
		} else if err != nil {
			return err
		}
	}

	return nil
}
