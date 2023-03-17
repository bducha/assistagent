//go:build linux
// +build linux

package system

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"syscall"
)

// Parse the output of hostnamectl to get the operating system
func getOS() (string, error) {
	

	// Execute command
	cmd := exec.Command("hostnamectl")
	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	// Parse output
	re := regexp.MustCompile(`(?i)static hostname: (?P<hostname>.*)[\S\s]*operating system: (?P<os>.*)\n`)

	matches := re.FindStringSubmatch(string(output))

	if len(matches) == 0 {
		return "", errors.New("no match found in hostnamectl output")
	}

	return matches[re.SubexpIndex("os")], nil
}

// Returns the total, used and free memory of the system
func GetMemory() (Memory, error) {
	info := &syscall.Sysinfo_t{}
	err := syscall.Sysinfo(info)
	if err != nil {
		return Memory{}, err
	}

	fmt.Println("total ram", info.Totalram)
	fmt.Println("free ram", info.Freeram)

	totalRam := uint64(info.Totalram) * uint64(info.Unit)
	freeRam := uint64(info.Freeram) * uint64(info.Unit)
	return Memory{
		TotalMemory: totalRam,
		FreeMemory: freeRam,
		UsedMemory: totalRam - freeRam,
	}, nil
}