package pkg

/*
#cgo LDFLAGS: -L./python312runtime -lpython312
#include <stdlib.h>
#include <stdio.h>
#include <Python.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func initPython() {
	C.Py_Initialize()
}

func shutdownPython() {
	C.Py_Finalize()
}

func runString(code string) error {
	cCode := C.CString(code)
	defer C.free(unsafe.Pointer(cCode))

	ret := C.PyRun_SimpleString(cCode)
	if ret != 0 {
		C.PyErr_Print()
		return fmt.Errorf("error running Python code")
	}
	return nil
}

func runFile(filename string) error {
	cFilename := C.CString(filename)
	defer C.free(unsafe.Pointer(cFilename))

	cMode := C.CString("r")
	defer C.free(unsafe.Pointer(cMode))

	file := C.fopen(cFilename, cMode)
	if file == nil {
		return fmt.Errorf("unable to open file: %s", filename)
	}
	defer C.fclose(file)

	ret := C.PyRun_SimpleFile(file, cFilename)
	if ret != 0 {
		C.PyErr_Print()
		return fmt.Errorf("error running Python file: %s", filename)
	}
	return nil
}
