package agent

import (
	"fmt"

	"github.com/bducha/assistagent/linux"
	"github.com/bducha/assistagent/system"
)

type Agent struct {
	system     system.System
	SystemInfo system.SystemInfo
}

func NewAgent() Agent {

	system := linux.LinuxSystem{}
	agent := Agent{
		system: &system,
	}

	// Init system infos
	info, err := agent.system.GetSysInfo()

	if err != nil {
		fmt.Println(err)
		panic("Error while trying to initialize system info. Stopping...")
	}

	agent.SystemInfo = info
	return agent
}

func (a *Agent) GetSysInfo() (system.SystemInfo, error) {
	return a.system.GetSysInfo()
}
