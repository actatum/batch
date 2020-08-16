package main

import (
	"fmt"
	"os"

	"github.com/actatum/batch/transport"
)

func main() {
	if err := transport.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
