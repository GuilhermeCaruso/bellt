package main

import (
	"fmt"

	"github.com/GuilhermeCaruso/bellt/cli/commands"
	"github.com/GuilhermeCaruso/bellt/cli/pkg"
)

var testCommandList = []pkg.Command{
	&commands.Generate{},
}

func main() {
	if err := pkg.StartCommandLine(testCommandList); err != nil {
		fmt.Println(err.Error())
	}
}
