package main

import (
	"fmt"

	_ "github.com/frijoles524/Blacklight/pkg"
)

func main() {
	initPython()
	defer shutdownPython()

	err := runString(`print("Hello from embedded Python!")`)
	if err != nil {
		fmt.Println("runString error:", err)
	}

	err = runFile("test.py")
	if err != nil {
		fmt.Println("runFile error:", err)
	}
}
