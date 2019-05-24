package commands

import (
	"flag"
	"fmt"
)

// Generate is a struct responsible for the document generation command line
type Generate struct {
	Type        int
	SourceFile  string
	OutputFile  string
	HelpCommand bool
}

const helpTextGenerate = `Generate is the command line to generate route documentation`
const helpLongTextGenerate = `
	Generate route documentation using a reference yaml file

	Generate documentation of:
		- routes
		- middlewares
		- static files

	bellt generate --t [documentationType] --f [documentationYamlFile] --o [outputNamefile]

	- documentationType: int
		0: all
		1: routesOnly
		2: middlewaresOnly
		3: staticFilesOnly

	- documentationYamlFile: string

	- outputNameFile: string
`

const exampleTextGenerate = `
	bellt generate --t 0 --f ./documentation.yaml --o ./documentation
	bellt generate --t 1 --f ./routesonly.yaml --o ./routesdocumentatio          
`

// Name is responsible for setting the name associated with
// the command line
func (cmd *Generate) Name() string { return "generate" }

// Example is responsible for setting the sample text associated
// with the command line
func (cmd *Generate) Example() string { return exampleTextGenerate }

// Help is responsible for setting the help text associated with
// the command line.
func (cmd *Generate) Help() string { return helpTextGenerate }

// LongHelp is responsible for setting the long help text associated with
// the command line.
func (cmd *Generate) LongHelp() string { return helpLongTextGenerate }

// Register is the inference action of flags to the system
func (cmd *Generate) Register(fs *flag.FlagSet) {
	fs.IntVar(&cmd.Type, "t", 0, "define the type of documentation")
	fs.StringVar(&cmd.SourceFile, "f", "./belltdoc.yaml", "define the source file location")
	fs.StringVar(&cmd.OutputFile, "o", "./documentation", "define the document destination")
	fs.BoolVar(&cmd.HelpCommand, "help", false, "show a help documentation")
}

// Run is the command line execution function
func (cmd *Generate) Run() {
	if cmd.HelpCommand {
		fmt.Println(cmd.LongHelp())
		return
	}

	if cmd.SourceFile == "" || cmd.OutputFile == "" {
		fmt.Println("Please inform the location of the documentation file and its output.")
		return
	}
}
