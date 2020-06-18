package pkg

import (
	"os"
	"testing"

	"github.com/GuilhermeCaruso/bellt/bellt/commands"
)

var emptyCommandList = []Command{}
var emptyOnlyOneCommand = []Command{
	&commands.Generate{},
}

func TestEmptyStartCommandLine(t *testing.T) {
	want := "Command line initialization require one or more commands"
	err := StartCommandLine(emptyCommandList)
	if err.Error() != want {
		t.Errorf("Error handling error: want %s, got %s", want, err)
	}
}

func TestStartCommandLineWithWrongArgumentLength(t *testing.T) {
	want := "Please pass some command"
	os.Args = []string{
		"--t=123",
	}
	err := StartCommandLine(emptyOnlyOneCommand)
	if err.Error() != want {
		t.Errorf("Error handling error: want %s, got %s", want, err.Error())
	}
}

func TestStartCommandLineWithRightArgumentLength(t *testing.T) {
	os.Args = []string{
		"bellt",
		"generate",
		"--t=123",
		"--f=456",
	}
	err := StartCommandLine(emptyOnlyOneCommand)
	if err != nil {
		t.Errorf("Error handling error: want %s, got %s", "nil", err.Error())
	}
}

func TestStartCommandLineWithoutCommands(t *testing.T) {
	want := "Please pass some command"
	os.Args = []string{
		"--test=123",
		"--test2=456",
	}
	err := StartCommandLine(emptyOnlyOneCommand)
	if err.Error() != want {
		t.Errorf("Error handling error: want %s, got %s", want, err.Error())
	}
}

func TestStartCommandLineWithCommand(t *testing.T) {
	os.Args = []string{
		"bellt",
		"generate",
		"--t=456",
	}
	err := StartCommandLine(emptyOnlyOneCommand)
	if err != nil {
		t.Errorf("Error handling error: want %s, got %s", "nil", err.Error())
	}
}

func TestStartCommandLineWithHelpCommand(t *testing.T){
	os.Args = []string{
		"bellt",
		"help",
	}

	err := StartCommandLine(emptyOnlyOneCommand)
	if err != nil {
		t.Errorf("Error handling error: want %s, got %s", "nil", err.Error())
	}
}

func TestArgumentFilter(t *testing.T) {
	want := UserCommand{
		Command: "generate",
		Arguments: []string{
			"-t=0",
		},
	}

	arguments := []string{
		"generate",
		"--t=0",
	}

	command := ArgumentsFilter(arguments)

	if command.Command != want.Command {
		t.Errorf("Command configuration error: want %s, got %s",
			want.Command,
			command.Command)
		return
	}

	if len(command.Arguments) == 0 {
		t.Errorf("Arguments configuration error")
		return
	}

	if want.Arguments[0] != command.Arguments[0] {
		t.Errorf("Arguments configuration error: want %s, got %s", want.Arguments[0], command.Arguments[0])
		return
	}
}
