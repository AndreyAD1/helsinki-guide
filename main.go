package main

import (
	"os"

	"github.com/AndreyAD1/helsinki-guide/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
