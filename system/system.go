package system

type System interface {
	GetSysInfo() (SystemInfo, error)
}

type SystemInfo struct {
	OS string
	Hostname string
}