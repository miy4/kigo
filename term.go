package kigo

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/sys/unix"
)

type Terminal struct {
	org  *unix.Termios
	in   *os.File
	out  *os.File
	size *WinSize
	buf  *strings.Builder
	keys chan Key
	errs chan error
}

type WinSize struct {
	unix.Winsize
}

type keyScanner struct {
	in   *bufio.Reader
	keys chan Key
	errs chan error
}

func (sc *keyScanner) scan() {
	for {
		b, err := sc.in.ReadByte()
		if err != nil && err != io.EOF {
			sc.errs <- err
			break
		}

		if b == 0x1b {
			if err = sc.scanTrailingSeq(); err != nil {
				sc.errs <- err
				break
			}
			continue
		}
		sc.keys <- Key(b)
	}
}

func (sc *keyScanner) scanTrailingSeq() error {
	b, err := sc.in.ReadByte()
	if err != nil && err != io.EOF {
		return err
	} else if b != '[' {
		sc.keys <- Key(0x1b)
		sc.keys <- Key(b)
		return nil
	}

	b, err = sc.in.ReadByte()
	if err != nil && err != io.EOF {
		return err
	}

	switch b {
	case 'A': // Up:    \x1b[A
		sc.keys <- KeyUp
	case 'B': // Down:  \x1b[B
		sc.keys <- KeyDown
	case 'C': // Right: \x1b[C
		sc.keys <- KeyRight
	case 'D': // Left:  \x1b[D
		sc.keys <- KeyLeft
	default:
		sc.keys <- Key(0x1b)
		sc.keys <- Key('[')
		sc.keys <- Key(b)
	}

	return nil
}

func NewTerminal() *Terminal {
	return &Terminal{
		in:  os.Stdin,
		out: os.Stdout,
		buf: &strings.Builder{},
	}
}

func (term *Terminal) getWinSize() (*WinSize, error) {
	ws, err := unix.IoctlGetWinsize(int(term.out.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return nil, err
	} else if ws.Col == 0 {
		r, c, err := term.getCursorPosition()
		if err != nil {
			return nil, err
		}
		return &WinSize{unix.Winsize{Row: r, Col: c}}, nil
	}

	return &WinSize{*ws}, nil
}

func (term *Terminal) getCursorPosition() (uint16, uint16, error) {
	_, err := term.out.WriteString("\x1b[999C\x1b[999B\x1b[6n")
	if err != nil {
		return 0, 0, err
	}

	b := make([]byte, 1)
	var buf bytes.Buffer
	for i := 0; i < 32; i++ {
		n, err := term.in.Read(b)
		if n != 1 || err != nil || b[0] == 'R' {
			break
		}
		buf.WriteByte(b[0])
	}

	seq := buf.Bytes()
	if len(seq) < 2 || seq[0] != '\x1b' || seq[1] != '[' {
		return 0, 0, fmt.Errorf("unexpected output: %s", seq)
	}

	var r, c uint16
	_, err = fmt.Sscanf(string(seq[2:]), "%d;%d", &r, &c)
	if err != nil {
		return 0, 0, fmt.Errorf("unexpected output: %s", seq)
	}

	return r, c, nil
}

func (term *Terminal) init() error {
	ws, err := term.getWinSize()
	if err != nil {
		return err
	}

	term.size = ws

	keys := make(chan Key, 1)
	errs := make(chan error, 1)
	term.keys = keys
	term.errs = errs

	scanner := &keyScanner{
		in:   bufio.NewReader(term.in),
		keys: keys,
		errs: errs,
	}
	go scanner.scan()

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
	select {
	case key := <-term.keys:
		return key, nil
	case err := <-term.errs:
		return Key(0), err
	}
}

func (term *Terminal) clearEntireScreen() {
	term.buf.WriteString("\x1b[2J")
}

func (term *Terminal) clearLineRight() {
	term.buf.WriteString("\x1b[K")
}

func (term *Terminal) moveCursor(row, col int) {
	fmt.Fprintf(term.buf, "\x1b[%d;%dH", row+1, col+1)
}

func (term *Terminal) moveCursorToHome() {
	term.moveCursor(0, 0)
}

func (term *Terminal) hideCursor() {
	term.buf.WriteString("\x1b[?25l")
}

func (term *Terminal) showCursor() {
	term.buf.WriteString("\x1b[?25h")
}

func (term *Terminal) writeString(s string) {
	term.buf.WriteString(s)
}

func (term *Terminal) flush() {
	term.out.WriteString(term.buf.String())
	term.buf.Reset()
}
