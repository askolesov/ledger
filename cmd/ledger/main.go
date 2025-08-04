package main

import (
	"go-finances/pkg/command"
	"os"
)

func main() {
	if err := command.GetRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
