package main

import (
	"io"
	"os"

	"golang.org/x/sys/unix"
)

func enableRawMode() {
	raw, _ := unix.IoctlGetTermios(int(os.Stdin.Fd()), unix.TCGETS)
	raw.Lflag &^= unix.ECHO
	unix.IoctlSetTermios(int(os.Stdin.Fd()), unix.TCSETSF, raw)
}

func main() {
	enableRawMode()

	b := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(b)
		if err == io.EOF || b[0] == 'q' {
			break
		}
	}

	os.Exit(0)
}
