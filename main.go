package main

import (
	"os"

	"github.com/mnowotnik/fzshell/cmd"
)

func main() {
	os.Exit(cmd.Run())
}
