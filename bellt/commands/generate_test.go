package commands

import (
	"bytes"
	"flag"
	"log"
	"os"
	"testing"
)

type FlagTest struct {
	Flag             *flag.Flag
	ExpectedName     string
	ExpectedUsage    string
	ExpectedDefValue interface{}
}

var blankGenerateCommand = Generate{}

var listOfGenerateCommands = []Generate{
	Generate{
		Type:       0,
		SourceFile: "./test/documentation/doc.yaml",
		OutputFile: "./test/documentation/output/new",
	},
	Generate{
		Type:       1,
		SourceFile: "./test/documentation/docTwo.yaml",
		OutputFile: "./test/documentation/output/newTwo",
	},
	Generate{
		Type:       2,
		SourceFile: "./test/documentation/docThree.yaml",
		OutputFile: "./test/documentation/output/newThree",
	},
	Generate{
		Type:       3,
		SourceFile: "./test/documentation/docFour.yaml",
		OutputFile: "./test/documentation/output/newFour",
	},
}

func TestCommandName(t *testing.T) {
	got := blankGenerateCommand.Name()
	want := "generate"

	if want != got {
		t.Errorf("command name is wrong: got %v want %v",
			got, want)
	}
}

func TestCommandHelpText(t *testing.T) {
	got := blankGenerateCommand.Help()
	want := `Generate is the command line to generate route documentation`

	if want != got {
		t.Errorf("command help text is wrong: got %v want %v",
			got, want)
	}

}

func TestCommandHelpLongText(t *testing.T) {
	got := blankGenerateCommand.LongHelp()
	want := `
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
	if want != got {
		t.Errorf("command help with long text is wrong: got %v want %v",
			got, want)
	}
}

func TestCommandExampleText(t *testing.T) {
	got := blankGenerateCommand.Example()
	want := `
	bellt generate --t 0 --f ./documentation.yaml --o ./documentation
	bellt generate --t 1 --f ./routesonly.yaml --o ./routesdocumentatio          
`

	if want != got {
		t.Errorf("command example is wrong: got %v want %v",
			got, want)
	}
}

func TestCommandLineCommands(t *testing.T) {
	flagSet := flag.NewFlagSet(blankGenerateCommand.Name(), flag.ContinueOnError)

	blankGenerateCommand.Register(flagSet)

	listOfFlags := []FlagTest{
		FlagTest{
			Flag:             flagSet.Lookup("t"),
			ExpectedName:     "t",
			ExpectedUsage:    "define the type of documentation",
			ExpectedDefValue: "0",
		},
		FlagTest{
			Flag:             flagSet.Lookup("f"),
			ExpectedName:     "f",
			ExpectedUsage:    "define the source file location",
			ExpectedDefValue: "./belltdoc.yaml",
		},
		FlagTest{
			Flag:             flagSet.Lookup("o"),
			ExpectedName:     "o",
			ExpectedUsage:    "define the document destination",
			ExpectedDefValue: "./documentation",
		},
	}

	for _, flag := range listOfFlags {
		if flag.Flag.Name != flag.ExpectedName {
			t.Errorf("wrong flag name: got %v want %v",
				flag.Flag.Name, flag.ExpectedName)
		}

		if flag.Flag.Usage != flag.ExpectedUsage {
			t.Errorf("wrong flag usage: got %v want %v",
				flag.Flag.Usage, flag.ExpectedUsage)
		}

		if flag.Flag.DefValue != flag.ExpectedDefValue {
			t.Errorf("wrong flag default value: got %v want %v",
				flag.Flag.DefValue, flag.ExpectedDefValue)
		}
	}
}

func TestRunCommand(t *testing.T) {
	flagSet := flag.NewFlagSet(blankGenerateCommand.Name(), flag.ContinueOnError)
	blankGenerateCommand.Register(flagSet)
	flagSet.Parse([]string{
		"--f",
		"",
		"--o",
		"",
	})

	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	blankGenerateCommand.Run()
	t.Log(buf.String())

}
