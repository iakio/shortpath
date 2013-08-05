package main

import (
	"flag"
	"fmt"
	"syscall"
	"unsafe"
)

var (
	kernel32         = syscall.MustLoadDLL("kernel32.dll")
	getShortPathName = kernel32.MustFindProc("GetShortPathNameW")
)

func abort(funcname string, err error) {
	panic(fmt.Sprintf("%s failed: %v", funcname, err))
}

func GetShortPathName(longpath string) string {
	var (
		bufferLen, ret uintptr
		lastErr        error
	)

	bufferLen, _, lastErr = getShortPathName.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(longpath))),
		0,
		0)

	if bufferLen == 0 {
		abort("GetShortPathName", lastErr)
	}
	shortpath := make([]uint16, bufferLen)

	ret, _, lastErr = getShortPathName.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(longpath))),
		uintptr(unsafe.Pointer(&shortpath[0])),
		bufferLen)

	if ret == 0 {
		abort("GetShortPathName", lastErr)
	}

	return syscall.UTF16ToString(shortpath)
}

func main() {
	defer kernel32.Release()

	flag.Parse()
	for _, longPath := range flag.Args() {
		fmt.Println(GetShortPathName(longPath))
	}
}
