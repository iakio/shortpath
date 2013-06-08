package main

import (
	"flag"
	"fmt"
	"syscall"
	"unsafe"
)

var (
	kernel32, _         = syscall.LoadLibrary("kernel32.dll")
	getShortPathName, _ = syscall.GetProcAddress(kernel32, "GetShortPathNameW")
)

func abort(funcname string, err error) {
	panic(fmt.Sprintf("%s failed: %v", funcname, err))
}

func GetShortPathName(longpath string) string {

	ret, _, callErr := syscall.Syscall(uintptr(getShortPathName),
		3,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(longpath))),
		0,
		0)

	if callErr != 0 {
		abort("GetShortPathName", callErr)
	}
	shortpath := make([]uint16, ret)

	ret, _, callErr = syscall.Syscall(uintptr(getShortPathName),
		3,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(longpath))),
		uintptr(unsafe.Pointer(&shortpath[0])),
		ret)

	if callErr != 0 {
		abort("GetShortPathName", callErr)
	}

	return syscall.UTF16ToString(shortpath)
}

func main() {
	defer syscall.FreeLibrary(kernel32)

	flag.Parse()
	for _, longPath := range flag.Args() {
		fmt.Println(GetShortPathName(longPath))
	}
}
