package main

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows/registry"
)

var (
	kernel32                   = syscall.NewLazyDLL("kernel32.dll")
	setConsoleScreenBufferSize = kernel32.NewProc("SetConsoleScreenBufferSize")
)

type coord struct {
	x int16
	y int16
}

func fatalIf(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}
}

func main() {
	k, err := registry.OpenKey(registry.CURRENT_USER, "Console", registry.QUERY_VALUE)
	fatalIf(err)
	defer k.Close()

	v, _, err := k.GetIntegerValue("ScreenBufferSize")
	fatalIf(err)
	col := int(uint16(v))
	row := int(uint16(v >> 16 & 0xffff))

	switch len(os.Args) {
	case 2:
		row, _ = strconv.Atoi(os.Args[1])
	case 3:
		row, _ = strconv.Atoi(os.Args[1])
		col, _ = strconv.Atoi(os.Args[2])
	}

	out, err := syscall.Open("CONOUT$", syscall.O_RDWR, 0)
	fatalIf(err)
	defer syscall.CloseHandle(out)

	size := coord{x: int16(col), y: int16(row)}
	r1, _, err := setConsoleScreenBufferSize.Call(uintptr(out), uintptr(*(*int32)(unsafe.Pointer(&size))))
	if r1 == 0 && err != nil {
		fatalIf(err)
	}
}
