package main

/*
#include <stdlib.h>
#include "test.h"

extern void startCgo(int);
extern void endCgo(int, int);
*/
import "C"
import (
	"fmt"
	"unsafe"
)

//export goStart
func goStart(i C.int) {
	fmt.Println(int(i))
}

//export goEnd
func goEnd(a C.int, b C.int) {
	fmt.Println(int(a), int(b))
}

func GoTraverse(filename string) {
	cCallbacks := C.Callbacks{}

	cCallbacks.start = C.StartCallbackFn(C.startCgo)
	cCallbacks.end = C.EndCallbackFn(C.endCgo)

	var cfilename *C.char = C.CString(filename)
	defer C.free(unsafe.Pointer(cfilename))

	C.traverse(cfilename, cCallbacks)
}

func main() {
	GoTraverse("test")
}
