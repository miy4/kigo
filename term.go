package kigo

import (
	"errors"
	"fmt"
	"io"
	"os"

	"golang.org/x/sys/unix"
)

type Terminal struct {
	org  *unix.Termios
	in   *os.File
	out  *os.File
	size *WinSize
}

type WinSize struct {
	unix.Winsize
}

func NewTerminal() *Terminal {
	return &Terminal{
		in:  os.Stdin,
		out: os.Stdout,
	}
}

func getWinSize(out *os.File) (*WinSize, error) {
	ws, err := unix.IoctlGetWinsize(int(out.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return nil, err
	} else if ws.Col == 0 {
		return nil, errors.New("possible errornous outcome")
	}

	return &WinSize{*ws}, nil
}

func (term *Terminal) init() error {
	ws, err := getWinSize(term.out)
	if err != nil {
		return err
	}

	term.size = ws
	return nil
}

func (term *Terminal) EnableRawMode() error {
	raw, err := unix.IoctlGetTermios(int(term.in.Fd()), unix.TCGETS)
	if err != nil {
		return err
	}

	org := *raw
	term.org = &org

	raw.Iflag &^= unix.BRKINT | unix.IXON | unix.INPCK | unix.ISTRIP | unix.ICRNL
	raw.Oflag &^= unix.OPOST
	raw.Cflag |= unix.CS8
	raw.Lflag &^= unix.ECHO | unix.ICANON | unix.IEXTEN | unix.ISIG
	raw.Cc[unix.VMIN] = 0
	raw.Cc[unix.VTIME] = 1

	err = unix.IoctlSetTermios(int(term.in.Fd()), unix.TCSETSF, raw)
	if err != nil {
		return err
	}

	return nil
}

func (term *Terminal) DisableRawMode() error {
	if term.org == nil {
		return nil
	}

	err := unix.IoctlSetTermios(int(term.in.Fd()), unix.TCSETSF, term.org)
	if err != nil {
		return err
	}

	return nil
}

func (term *Terminal) ReadKey() (Key, error) {
	b := make([]byte, 1)
	for {
		n, err := term.in.Read(b)
		if err != nil && err != io.EOF {
			return 0, err
		} else if n == 1 {
			break
		}
	}
	return Key(b[0]), nil
}

func (term *Terminal) clearEntireScreen() {
	term.out.WriteString("\x1b[2J")
}

func (term *Terminal) moveCursor(row, col int) {
	fmt.Fprintf(term.out, "\x1b[%d;%dH", row+1, col+1)
}

func (term *Terminal) moveCursorToHome() {
	term.moveCursor(0, 0)
}

func (term *Terminal) writeString(s string) {
	term.out.WriteString(s)
}