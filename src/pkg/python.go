package main

/*
#cgo LDFLAGS: -L./python312runtime -lpython312
#include <stdlib.h>
#include <stdio.h>

extern void Py_Initialize(void);
extern void Py_Finalize(void);
extern int PyRun_SimpleString(const char *command);
extern int PyRun_SimpleFile(FILE *fp, const char *filename);
extern void PyErr_Print(void);
*/
import "C"

import (
	"fmt"
	"unsafe"
)

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

	file := C.fopen(cFilename, C.CString("r"))
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
