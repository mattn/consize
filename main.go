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
	setConsoleWindowInfo       = kernel32.NewProc("SetConsoleWindowInfo")
)

type COORD struct {
	X int16
	Y int16
}

type SMALL_RECT struct {
	Left   int16
	Top    int16
	Right  int16
	Bottom int16
}

func fatalIf(f string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s: %v\n", os.Args[0], f, err)
		os.Exit(1)
	}
}

func main() {
	k, err := registry.OpenKey(registry.CURRENT_USER, "Console", registry.QUERY_VALUE)
	fatalIf("RegOpenKeyEx", err)
	defer k.Close()

	v, _, err := k.GetIntegerValue("ScreenBufferSize")
	fatalIf("RegQueryValueEx", err)
	col := int(uint16(v))
	row := int(uint16(v >> 16 & 0xffff))

	switch len(os.Args) {
	case 2:
		row, _ = strconv.Atoi(os.Args[1])
	case 3:
		row, _ = strconv.Atoi(os.Args[1])
		col, _ = strconv.Atoi(os.Args[2])
	}

	out, err := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	if err != nil {
		out, err = syscall.Open("CONOUT$", syscall.O_RDWR, 0)
		fatalIf("CreateFile", err)
	}
	defer syscall.CloseHandle(out)

	rect := SMALL_RECT{Left: int16(1), Top: int16(1), Right: int16(col), Bottom: int16(25)}
	r1, _, err := setConsoleWindowInfo.Call(uintptr(out), uintptr(int32(1)), uintptr(unsafe.Pointer(&rect)))
	if r1 == 0 && err != nil {
		fatalIf("SetConsoleWindowInfo", err)
	}

	size := COORD{X: int16(col), Y: int16(row)}
	r1, _, err = setConsoleScreenBufferSize.Call(uintptr(out), *(*uintptr)(unsafe.Pointer(&size)))
	if r1 == 0 && err != nil {
		fatalIf("SetConsoleScreenBufferSize", err)
	}
}
