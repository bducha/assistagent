//go:build windows
// +build windows

package system

import (
	"errors"
	"os/exec"
	"regexp"
)

// Parse the output of the "ver" command on windows to get the OS
func getOS() (string, error) {
	cmd := exec.Command("cmd", "ver")

	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	// Parse output
	re := regexp.MustCompile(`(?i)^(?P<os>.*) \[version (?P<version>\d*)`)

	matches := re.FindStringSubmatch(string(output))

	if len(matches) == 0 {
		return "", errors.New("no match found in hostnamectl output")
	}

	return matches[re.SubexpIndex("os")] + " " + matches[re.SubexpIndex("version")], nil
}