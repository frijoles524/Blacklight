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

	version := *version
	if version == "" {
		var err error
		version, err = pkg.GetHighestVersion(*software)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	if !pkg.IsInstalled(&pkg.Package{Name: *software, Version: version}) {
		fmt.Printf("Software %s version %s is not installed.\n", *software, version)
		return
	}

	err := pkg.RunFile(fmt.Sprintf("%s-%s", *software, version))
	if err != nil {
		fmt.Printf("Error running software: %v\n", err)
	}
}
