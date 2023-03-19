package main

import (
	"fmt"
	"os"

	"github.com/miy4/kigo"
)

func main() {
	err := kigo.NewEditor().Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running the editor: %v\n", err)
		os.Exit(1)
	}
}
