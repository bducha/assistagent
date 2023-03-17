//go:build linux
// +build linux

package system

import (
	"errors"
	"os/exec"
	"regexp"
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