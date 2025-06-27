package main

import (
	"fmt"

	"github.com/frijoles524/Blacklight/pkg"
)

func RunApp() {
	//pkg.InitPython()
	defer pkg.ShutdownPython()

	if software == nil || *software == "" {
		fmt.Println("No software name was passed.")
		return
	}

	var versionStr string
	if version == nil || *version == "" {
		var err error
		versionStr, err = pkg.GetHighestVersion(*software)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		versionStr = *version
	}

	if !pkg.IsInstalled(&pkg.Package{Name: *software, Version: versionStr}) {
		fmt.Printf("Software %s version %s is not installed.\n", *software, versionStr)
		return
	}

	<-pythonInitDone
	err := pkg.RunFile(fmt.Sprintf("%s-%s/%s.py", *software, versionStr, *software))
	if err != nil {
		fmt.Printf("Error running software: %v\n", err)
	}
}
