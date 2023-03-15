package main

import (
	"fmt"

	"github.com/bducha/assistagent/agent"
)

func main() {

	agent := agent.NewAgent()

	info, _ := agent.GetSysInfo()

	fmt.Printf("System : %s, %s", info.Hostname, info.OS)
}