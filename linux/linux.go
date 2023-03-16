package linux

import (
	"errors"
	"os/exec"
	"regexp"

	"github.com/bducha/assistagent/system"
)

type LinuxSystem struct{}

func (s *LinuxSystem) GetSysInfo() (system.SystemInfo, error) {

	sysInfo := system.SystemInfo{}

	// Execute command
	cmd := exec.Command("hostnamectl")
	output, err := cmd.Output()

	if err != nil {
		return sysInfo, err
	}

	// Parse output
	re := regexp.MustCompile(`(?i)static hostname: (?P<hostname>.*)[\S\s]*operating system: (?P<os>.*)\n`)

	matches := re.FindStringSubmatch(string(output))

	if len(matches) == 0 {
		return sysInfo, errors.New("no match found in hostnamectl output")
	}

	os := matches[re.SubexpIndex("os")]
	hostname := matches[re.SubexpIndex("hostname")]

	return system.SystemInfo{OS: os, Hostname: hostname}, nil
}
