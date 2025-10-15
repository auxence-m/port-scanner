package main

import (
	"context"
	"log"
	"pScan/cmd"

	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra/doc"
)

func main() {
	//cmd.Execute()

	rootCmd := cmd.Root()

	if err := fang.Execute(context.Background(), rootCmd, fang.WithoutManpage(), fang.WithVersion(rootCmd.Version)); err != nil {
		log.Fatal(err)
	}

	err := doc.GenMarkdownTree(rootCmd, "./docs")
	if err != nil {
		log.Fatal(err)
	}
}
