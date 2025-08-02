package main

import "go-finances/pkg/command"

func main() {
	_ = command.GetRootCmd().Execute()
}
