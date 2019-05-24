// Copyright 2019 Guilherme Caruso. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package pkg

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"text/tabwriter"
)

// Command is an interface responsible for setting the  default format for the
// command line
type Command interface {
	Name() string
	Example() string
	Help() string
	LongHelp() string
	Register(*flag.FlagSet)
	Run()
}

// UserCommand is a standardized structure for
// command line calls
type UserCommand struct {
	Command   string
	Arguments []string
}

// StartCommandLine executes and initializes the entire command line
// routine of the application
func StartCommandLine(commands []Command) error {
	if len(commands) == 0 {
		return errors.
			New("Command line initialization require one or more commands")
	}

	if len(os.Args) < 2 {
		InitCommandText(commands)
		return errors.New("Please pass some command")
	}

	userPassedArguments := os.Args[1:]

	userCommand := ArgumentsFilter(userPassedArguments)

	if userCommand.Command == "" {
		return errors.New("Please pass some command")
	}
	for _, command := range commands {
		if userCommand.Command == "help" {
			InitCommandText(commands)
		}
		if userCommand.Command == command.Name() {
			flagSet := flag.NewFlagSet(command.Name(), flag.ContinueOnError)
			command.Register(flagSet)
			flagSet.Parse(os.Args[2:])
			command.Run()
		}
	}

	return nil
}

// ArgumentsFilter is responsible for standardizing the call of users
func ArgumentsFilter(commands []string) UserCommand {
	regularExpressionValidator := regexp.MustCompile("(?m)\\-")
	commandSet := false

	var commandDefinition UserCommand

	for _, argument := range commands {
		if !regularExpressionValidator.MatchString(argument) && !commandSet {
			commandDefinition.Command = argument
			commandSet = true
		}

		if regularExpressionValidator.MatchString(argument) {
			commandDefinition.Arguments =
				append(commandDefinition.Arguments, argument)
		}
	}

	return commandDefinition
}

// InitCommandText returns a help text for the user at the
// command line.
func InitCommandText(commands []Command) {
	fmt.Println("Bellt")
	fmt.Println("Please use any of the listed commands:")
	tabWriter := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, command := range commands {
		fmt.Fprintf(tabWriter, "\t- %s\t%s\n", command.Name(), command.Help())
	}
	tabWriter.Flush()
	fmt.Println(``)

	for _, command := range commands {
		fmt.Fprintf(tabWriter, "\t%s\n", command.Example())
	}
}
