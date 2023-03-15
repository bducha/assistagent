package agent

import (
	"github.com/bducha/assistagent/system"
	"github.com/bducha/assistagent/linux"
)

type Agent struct {
	system system.System
}

func NewAgent() Agent {

	system := linux.LinuxSystem{}
	agent := Agent{
		system: &system,
	}
	return agent
}

func (a *Agent) GetSysInfo() (system.SystemInfo, error) {
	return a.system.GetSysInfo()
}
