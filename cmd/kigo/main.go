package main

import (
	"fmt"
	"os"
	"unicode"

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

func main() {
	term := NewTerminal()
	err := term.EnableRawMode()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error enabling raw mode: %v\n", err)
		os.Exit(1)
	}
	defer term.DisableRawMode()

	b := make([]byte, 1)
	for {
		b[0] = 0
		term.in.Read(b)

		r := rune(b[0])
		if unicode.IsPrint(r) {
			fmt.Printf("%d ('%c')\r\n", r, r)
		} else {
			fmt.Printf("%d\r\n", r)
		}

		if b[0] == 'q' {
			break
		}
	}
}
