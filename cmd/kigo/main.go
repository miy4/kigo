package main

import (
	"io"
	"os"
)

func main() {
	b := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(b)
		if err == io.EOF || b[0] == 'q' {
			break
		}
	}

	os.Exit(0)
}
