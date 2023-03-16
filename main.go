package main

import (
	"fmt"
	"time"

	"github.com/bducha/assistagent/agent"
	"github.com/bducha/assistagent/mqtt"
)

// Dev consts
const (
	OBJECT_ID = "pc-fixe"
)

func main() {

	agent := agent.NewAgent()

	fmt.Println("AssistAgent started")

	fmt.Printf("System : %s, %s", agent.SystemInfo.Hostname, agent.SystemInfo.OS)
	fmt.Println()

	mqtt := mqtt.NewClient()

	discoveryPayload := mqtt.GetBaseDiscoveryPayload()
	discoveryPayload.Name = "Hostname"
	discoveryPayload.UniqueId = mqtt.GetClientId() + "_" + "hostname"
	discoveryPayload.StateTopic = mqtt.GetStateTopic() + discoveryPayload.UniqueId

	if err := mqtt.PublishDiscoveryPayload(discoveryPayload, "text"); err != nil {
		fmt.Println(err)
	}

	time.Sleep(5 * time.Second)
	
	mqtt.PublishStateUpdate(discoveryPayload.UniqueId, agent.SystemInfo.Hostname)
}