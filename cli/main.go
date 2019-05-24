package main

import (
	"os"

	"github.com/GuilhermeCaruso/bellt/cli/commands"
	"github.com/GuilhermeCaruso/bellt/cli/pkg"
)

var testCommandList = []pkg.Command{
	&commands.Generate{},
}

func main() {
	if err := pkg.StartCommandLine(testCommandList); err != nil {
		os.Exit(1)
	}
}
