package main

import (
	"flag"
	"sgirl/internal/cli"
)

func main() {
	cli.Execute(flag.Args())
}
