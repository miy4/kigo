package kigo

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/sys/unix"
)

type Terminal struct {
	org *unix.Termios
	in  *os.File
	out *os.File
}

func NewTerminal() *Terminal {
	return &Terminal{
		in:  os.Stdin,
		out: os.Stdout,
	}
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
