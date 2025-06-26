package main

import (
	"flag"
)

var (
	command  = flag.String("command", "", "command to run")
	software = flag.String("software", "", "name of the software to perform the action on. not needed for every command")
	version  = flag.String("version", "", "version of the software in question")
)

func ParseFlags() {
	flag.Parse()
}
