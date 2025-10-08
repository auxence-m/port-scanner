package main

import (
	"log"
	"pScan/cmd"

	"github.com/spf13/cobra/doc"
)

func main() {
	cmd.Execute()

	rootCmd := cmd.Root()
	err := doc.GenMarkdownTree(rootCmd, "./docs")
	if err != nil {
		log.Fatal(err)
	}
}
