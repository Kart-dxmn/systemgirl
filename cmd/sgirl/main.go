package main

import (
	"os"
	"sgirl/internal/cli"
)

func main() {
	if len(os.Args) < 1 {
		cli.Help()
		return
	}

	switch os.Args[1] {
	case "add":
		if len(os.Args) < 6 {
			cli.Help()
		}
		cli.AddConnect()
	case "connect":
		if len(os.Args) < 3 {
			cli.Help()
		}
		cli.Connect()
	default:
		cli.Help()
	}
}
