package main

import (
	"github.com/actatum/batch/batch"
)

func main() {
	/*
		if err := transport.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	*/
	s, _ := batch.NewLogger()
	s.Info("testing zap logger")
	s.Error("this is an error")
	s.Fatal("am ded")
}
