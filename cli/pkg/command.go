// Copyright 2019 Guilherme Caruso. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.
package pkg

import (
	"flag"
)

// Command is an interface responsible for setting the  default format for the
// command line
type Command interface {
	Name() string
	Register(*flag.FlagSet)
	Help() string
	LongHelp() string
	Run()
	Example()
}
