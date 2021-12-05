package main

// #include <stdlib.h>
// #include "linuxcnc.hh"
// #cgo CPPFLAGS: -I${SRCDIR}/../../adapter
// #cgo LDFLAGS: -L${SRCDIR}/../../lib -llinuxcncadapter
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	statHandle := C.stat_newHandle()
	defer C.free(unsafe.Pointer(statHandle))

	if r := C.stat_initHandle(statHandle); r != 0 {
		panic("error during handle init")
	}

	if r := C.stat_poll(statHandle); r != 0 {
		panic("error calling stat_poll()")
	}

	// Get stuff
	global := C.struct_stat_Global{}
	C.stats_global(statHandle, &global)
	fmt.Printf("global:\n %#v\n\n", global)

	joints := make([]C.struct_stat_Joint, global.numberOfJoints)
	C.stats_joints(statHandle, &joints[0])
	for i, joint := range joints {
		fmt.Printf("joint #%d %#v\n", i, joint)
	}
	fmt.Println()

	spindles := make([]C.struct_stat_Spindle, global.numberOfSpindles)
	C.stats_spindles(statHandle, &spindles[0])
	for i, spindle := range spindles {
		fmt.Printf("spindle #%d %#v\n", i, spindle)
	}
}
