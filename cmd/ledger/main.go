package main

import (
	"ledger/pkg/command"
	"os"
)

func main() {
	if err := command.GetRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
