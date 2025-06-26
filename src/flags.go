package main

import (
	"flag"
)

var (
	command  = flag.String("command", "help", "command to run, leave empty for help")
	software = flag.String("software", "", "name of the software to perform the action on. not needed for every command")
)

func ParseFlags() {
	flag.Parse()
}
