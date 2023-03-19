package kigo

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

const version = "0.1.0"

var Quit = errors.New("Quit")

type row struct {
	chars string
}

type document struct {
	rows []*row
}

func newDocument() *document {
	return &document{
		rows: []*row{},
	}
}

func (doc *document) isEmpty() bool {
	return len(doc.rows) == 0
}

type Editor struct {
	term     *Terminal
	cur      *pos
	doc      *document
	FileName string
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
		doc:  newDocument(),
	}
}

func (editor *Editor) open(name string) error {
	f, err := os.Open(name)
	if err != nil {
		return err
	}

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		editor.doc.rows = append(editor.doc.rows, &row{chars: sc.Text()})
	}
	if err = sc.Err(); err != nil {
		return err
	}

	return nil
}

func (editor *Editor) rows() int {
	return int(editor.term.size.Row)
}

func (editor *Editor) cols() int {
	return int(editor.term.size.Col)
}

func (editor *Editor) drawRow(at int) {
	row := editor.doc.rows[at]
	width := editor.cols()
	if len(row.chars) < width {
		width = len(row.chars)
	}

	editor.term.writeString(row.chars[:width])
}

func (editor *Editor) drawRows() {
	for y := 0; y < editor.rows(); y++ {
		if y < len(editor.doc.rows) {
			editor.drawRow(y)
		} else {
			if editor.doc.isEmpty() && y == editor.rows()/3 {
				editor.drawWelcome()
			} else {
				editor.term.writeString("~")
			}
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
	case KeyUp, KeyDown, KeyRight, KeyLeft:
		editor.moveCursor(key)
	case KeyPgUp, KeyPgDn:
		dir := KeyDown
		if key == KeyPgUp {
			dir = KeyUp
		}
		for i := editor.rows(); i > 0; i-- {
			editor.moveCursor(dir)
		}
	case KeyHome:
		editor.cur.x = 0
	case KeyEnd:
		editor.cur.x = editor.cols() - 1
	}

	return nil
}

func (editor *Editor) moveCursor(key Key) {
	switch key {
	case KeyLeft:
		if editor.cur.x >= 1 {
			editor.cur.x--
		}
	case KeyRight:
		if editor.cur.x <= editor.cols()-2 {
			editor.cur.x++
		}
	case KeyUp:
		if editor.cur.y >= 1 {
			editor.cur.y--
		}
	case KeyDown:
		if editor.cur.y <= editor.rows()-2 {
			editor.cur.y++
		}
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

	if editor.FileName != "" {
		err := editor.open(editor.FileName)
		if err != nil {
			return err
		}
	}

	for {
		editor.refreshScreen()
		if err := editor.processKeypress(); err != nil {
			if err == Quit {
				break
			}

			return err
		}
	}

	return nil
}
