package agent

import (
	"fmt"

	"github.com/bducha/assistagent/system"
)

type Agent struct {
	SystemInfo system.SystemInfo
}

func NewAgent() Agent {


	agent := Agent{}

	// Init system infos
	info, err := system.GetSysInfo()

	if err != nil {
		fmt.Println(err)
		panic("Error while trying to initialize system info. Stopping...")
	}

	agent.SystemInfo = info
	return agent
}

func (a *Agent) GetSysInfo() system.SystemInfo {
	return a.SystemInfo
}
