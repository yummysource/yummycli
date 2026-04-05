package main

import (
	"os"

	"github.com/yummysource/yummycli/cmd"
)

func main() {
	os.Exit(cmd.Execute())
}
