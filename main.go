package main

import (
	"os"

	"github.com/sol-eng/wbi/cmd"
)

// https://goreleaser.com/cookbooks/using-main.version
var (
	version string = "dev"
	commit  string = "none"
	date    string = "unknown"
)

func main() {
	cmd.Execute(version, os.Args[1:])
}
