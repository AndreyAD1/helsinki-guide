package main

import (
	"fmt"
	"os"

	"github.com/AndreyAD1/helsinki-guide/cmd"
)


func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}