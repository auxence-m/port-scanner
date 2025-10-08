package main

import (
	"context"
	"log"
	"os"
	"pScan/cmd"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra/doc"
)

func main() {
	//cmd.Execute()

	rootCmd := cmd.Root()
	err := doc.GenMarkdownTree(rootCmd, "./docs")
	if err != nil {
		log.Fatal(err)
	}

	if err := fang.Execute(context.Background(), rootCmd); err != nil {
		os.Exit(1)
	}
}
