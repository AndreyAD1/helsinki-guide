package main

import (
	"log"
	"os"

	"github.com/spf13/cobra/doc"

	"github.com/AndreyAD1/helsinki-guide/cmd"
)

func main() {
	docsDirPath := "./command_docs"
	os.Mkdir(docsDirPath, 0750)
	err := doc.GenMarkdownTree(cmd.RootCmd, docsDirPath)
	if err != nil {
		log.Fatal(err)
	}
}