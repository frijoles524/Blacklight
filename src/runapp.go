package main

import (
	"fmt"

	"github.com/frijoles524/Blacklight/pkg"
)

func RunApp() {
	pkg.InitPython()
	defer pkg.ShutdownPython()
	if *software == "" {
		fmt.Println("No software name was passed.")
		return
	}
}
