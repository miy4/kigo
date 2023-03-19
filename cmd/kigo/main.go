package main

import (
	"fmt"
	"os"

	"github.com/miy4/kigo"
)

func main() {
	editor := kigo.NewEditor()
	if len(os.Args) >= 2 {
		editor.FileName = os.Args[1]
	}

	err := editor.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running the editor: %v\n", err)
		os.Exit(1)
	}
}
