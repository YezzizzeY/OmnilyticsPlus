package main

/*
#cgo CXXFLAGS: -std=c++11
#cgo LDFLAGS: -lntl -lgmp -lm -pthread
#include "interpolate.h"
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

func main() {
	prime := "23"
	xVals := []string{"1", "3", "4"}
	yVals := []string{"7", "6", "0"}

	primeC := C.CString(prime)
	defer C.free(unsafe.Pointer(primeC))

	xValsC := make([]*C.char, len(xVals))
	for i, v := range xVals {
		xValsC[i] = C.CString(v)
		defer C.free(unsafe.Pointer(xValsC[i]))
	}

	yValsC := make([]*C.char, len(yVals))
	for i, v := range yVals {
		yValsC[i] = C.CString(v)
		defer C.free(unsafe.Pointer(yValsC[i]))
	}

	result := C.interpolate(primeC, C.int(len(xVals)), &xValsC[0], &yValsC[0])
	fmt.Println(C.GoString(result))
}
