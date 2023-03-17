package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/bducha/assistagent/mqtt"
	"github.com/bducha/assistagent/system"
)

type Agent struct {
	SystemInfo system.SystemInfo
	mqtt mqtt.Client
}

func NewAgent() Agent {

	agent := Agent{}

	agent.mqtt = mqtt.NewClient()

	// Init system infos
	info, err := system.GetSysInfo()

	if err != nil {
		fmt.Println(err)
		panic("Error while trying to initialize system info. Stopping...")
	}

	agent.SystemInfo = info
	return agent
}

// Start the agent loop
func (a *Agent) Start(ctx context.Context) {
	
	a.discovery()
	a.mqtt.PublishStateUpdate(a.mqtt.GetClientId() + "_" + "hostname", a.SystemInfo.Hostname)
	
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// keep the agent available
			a.updateAvailability(true)
		case <-ctx.Done():
			// make the agent unavailable
			a.updateAvailability(false)
			return
		}

	}
}

// Send all discovery payloads for home assistant
func (a *Agent) discovery() {
	discoveryPayload := a.mqtt.GetBaseDiscoveryPayload()
	discoveryPayload.Name = "Hostname"
	discoveryPayload.UniqueId = a.mqtt.GetClientId() + "_" + "hostname"
	discoveryPayload.StateTopic = a.mqtt.GetStateTopic() + discoveryPayload.UniqueId

	if err := a.mqtt.PublishDiscoveryPayload(discoveryPayload, "sensor"); err != nil {
		fmt.Println(err)
	}
}

// updates the agent availability
func (a *Agent) updateAvailability(available bool) {
	a.mqtt.PublishAvailabilityPayload(available)
}


func (a *Agent) GetSysInfo() system.SystemInfo {
	return a.SystemInfo
}


