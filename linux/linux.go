package linux

import "github.com/bducha/assistagent/system"

type LinuxSystem struct {}

func (s *LinuxSystem) GetSysInfo() (system.SystemInfo, error) {
	return system.SystemInfo{OS: "Ubuntu 22.04", Hostname: "ben-desktop"}, nil
}