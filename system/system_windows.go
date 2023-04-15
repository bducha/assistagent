//go:build windows
// +build windows

package system

import (
	"fmt"
	"os/exec"
	"syscall"
	"unsafe"
)

type memStatusEx struct {
	dwLength     uint32
	dwMemoryLoad uint32
	ullTotalPhys uint64
	ullAvailPhys uint64
	unused       [5]uint64
}

func GetMemory() (Memory, error) {
    kernel32, err := syscall.LoadLibrary("kernel32.dll")
    if err != nil {
        panic(err)
    }
    defer syscall.FreeLibrary(kernel32)

    var memInfo memStatusEx
    memInfo.dwLength = uint32(unsafe.Sizeof(memInfo))

    ret, _, err := syscall.NewLazyDLL("kernel32.dll").NewProc("GlobalMemoryStatusEx").Call(uintptr(unsafe.Pointer(&memInfo)))
    if ret == 0 {
        panic(err)
    }

    totalMemory := memInfo.ullTotalPhys
    freeMemory := memInfo.ullAvailPhys
    usedMemory := totalMemory - freeMemory

    return Memory{TotalMemory: totalMemory, FreeMemory: freeMemory, UsedMemory: usedMemory}, nil
}

func Shutdown() {
	cmd := exec.Command("shutdown", "/s", "/t", "0")
    err := cmd.Run()
    if err != nil {
        fmt.Println("Error shutting down:", err)
    }
}