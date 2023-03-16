package main

import (
	"fmt"

	"github.com/bducha/assistagent/agent"
)

func main() {

	agent := agent.NewAgent()

	fmt.Println("AssistAgent started")

	fmt.Printf("System : %s, %s", agent.SystemInfo.Hostname, agent.SystemInfo.OS)
	fmt.Println()
}