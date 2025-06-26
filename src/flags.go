package main

import (
	"fmt"
	"os"
)

var (
	command  *string
	software *string
	version  *string
)

func PrintUsage() {
	fmt.Println(`Usage:
  blacklight.exe command [software] [version]

Arguments:
  command    Command to run (required)
  software   Name of the software (optional)
  version    Version of the software (optional)

Example:
  blacklight.exe install scratchpad`)
}

func ParseArgs() error {
	args := os.Args[1:]

	if len(args) == 0 {
		PrintUsage()
		os.Exit(0)
	}

	for _, arg := range args {
		if arg == "--help" || arg == "-h" {
			PrintUsage()
			os.Exit(0)
		}
	}

	command = &args[0]

	if len(args) > 1 {
		software = &args[1]
	} else {
		software = nil
	}
	if len(args) > 2 {
		version = &args[2]
	} else {
		version = nil
	}

	return nil
}
