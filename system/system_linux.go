//go:build linux
// +build linux

package system

import (
	"fmt"
	"os/exec"
	"syscall"
)

// Returns the total, used and free memory of the system
func GetMemory() (Memory, error) {
	info := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(info)
	if err != nil {
		return Memory{}, err
	}

	totalRam := uint64(info.Totalram) * uint64(info.Unit)
	freeRam := uint64(info.Freeram) * uint64(info.Unit)
	return Memory{
		TotalMemory: totalRam,
		FreeMemory: freeRam,
		UsedMemory: totalRam - freeRam,
	}, nil
}

// Shuts down the PC
func Shutdown() {
	cmd := exec.Command("shutdown", "now")
	err := cmd.Run()
	if err != nil {
        fmt.Println("Error shutting down:", err)
    }
}