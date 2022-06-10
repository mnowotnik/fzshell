package main

import (
	"os"

	"github.com/mnowotnik/fzshell/cmd"
)

var version string = "devel"
var revision string = "devel"

func main() {
	os.Exit(cmd.Run(version, revision))
}
