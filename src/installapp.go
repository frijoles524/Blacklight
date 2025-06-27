package main

import (
	"fmt"

	"github.com/frijoles524/Blacklight/pkg"
)

func InstallApp() {
	if software == nil || *software == "" {
		fmt.Println("No software name was passed.")
		return
	}

	p := &pkg.Package{
		Name:    *software,
		Version: "",
	}
	if version != nil {
		p.Version = *version
	}

	<-pythonInitDone
	if err := pkg.Download(p); err != nil {
		fmt.Println("Failed to install:", err)
	}
}
