package system

import "os"

type SystemInfo struct {
	OS       string
	Hostname string
}

// Get the system infos
func GetSysInfo() (SystemInfo, error) {

	hostname, err := os.Hostname()
	if err != nil {
		return SystemInfo{}, err
	}

	os, err := getOS()
	if err != nil {
		return SystemInfo{}, err
	}

	return SystemInfo{
		OS:       os,
		Hostname: hostname,
	}, nil
}