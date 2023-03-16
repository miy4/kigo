package main

import (
	"io"
	"os"
)

func main() {
	b := make([]byte, 1)
	for {
		_, err := os.Stdin.Read(b)
		if err == io.EOF {
			break
		}
	}

	os.Exit(0)
}
