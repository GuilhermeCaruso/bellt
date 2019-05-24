package main

import (
	"os"

	"github.com/GuilhermeCaruso/bellt/bellt/commands"
	"github.com/GuilhermeCaruso/bellt/bellt/pkg"
)

var testCommandList = []pkg.Command{
	&commands.Generate{},
}

func main() {
	if err := pkg.StartCommandLine(testCommandList); err != nil {
		os.Exit(1)
	}
}
