package agent

import (
	"context"
	"fmt"
	"time"

	"github.com/bducha/assistagent/mqtt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Agent struct {
	mqtt mqtt.Client
	sensorCollectors []SensorCollector
	buttons []Button
}

type SensorCollector struct {
	Name string
	UniqueId string
	DeviceClass string
	UnitOfMeasurement string
	SensorAttributes []SensorAttribute
	// The function that will be called to get the state payload
	// If there is multiple attributes, the function must return a json string with 
	// values for each attributes with the unique id as the key
	CollectorFunc func() (string , error)
}

type SensorAttribute struct {
	DisplayName string
	UniqueId string
	DeviceClass string
	UnitOfMeasurement string
}

type Button struct {
	Name string
	UniqueId string
	Action func(MQTT.Client, MQTT.Message)
}

func NewAgent() Agent {

	agent := Agent{}

	agent.mqtt = mqtt.NewClient()

	return agent
}

// Start the agent loop
func (a *Agent) Start(ctx context.Context) {
	
	a.sensorCollectors = []SensorCollector{
		{
			Name: "Memory",
			UniqueId: "memory",
			CollectorFunc: GetMemoryState,
			SensorAttributes: []SensorAttribute{
				{
					DisplayName: "Total memory",
					UniqueId: "total_memory",
					DeviceClass: "data_size",
					UnitOfMeasurement: "bit",
				},
				{
					DisplayName: "Free memory",
					UniqueId: "free_memory",
					DeviceClass: "data_size",
					UnitOfMeasurement: "bit",
				},
				{
					DisplayName: "Used memory",
					UniqueId: "used_memory",
					DeviceClass: "data_size",
					UnitOfMeasurement: "bit",
				},
			},
		},
	}

	a.buttons = []Button{
		{
			Name: "Shutdown",
			UniqueId: "shutdown",
			Action: Shutdown,
		},
	}

	a.registerSensorCollectors()
	a.registerButtons()
	
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// keep the agent available
			a.updateAvailability(true)
			a.updateSensorsState()
		case <-ctx.Done():
			// make the agent unavailable
			a.updateAvailability(false)
			return
		}

	}
}

// The function send discovery payload for every state collectors
func (a *Agent) registerSensorCollectors() {
	for _, el := range a.sensorCollectors {
		payload := a.mqtt.GetBaseDiscoveryPayload()
		payload.Name = el.Name
		payload.UniqueId = a.mqtt.GetClientId() + "_" + el.UniqueId
		payload.StateTopic = a.mqtt.GetStateTopic() + el.UniqueId
		payload.DeviceClass = el.DeviceClass
		payload.UnitOfMeasurement = el.UnitOfMeasurement

		if el.SensorAttributes != nil {
			for _, attr := range el.SensorAttributes {
				payload.Name = attr.DisplayName
				payload.UniqueId = attr.UniqueId
				payload.DeviceClass = attr.DeviceClass
				payload.UnitOfMeasurement = attr.UnitOfMeasurement
				payload.ValueTemplate = "{{ value_json." + attr.UniqueId + " }}"

				if err := a.mqtt.PublishDiscoveryPayload(payload, "sensor"); err != nil {
					fmt.Println(err)
				}
			}
			continue
		}
		if err := a.mqtt.PublishDiscoveryPayload(payload, "sensor"); err != nil {
			fmt.Println(err)
		}
	}
}

// Sends discovery payloads for each buttons configured
func (a *Agent) registerButtons() {
	for _, el := range a.buttons {
		payload := a.mqtt.GetBaseDiscoveryPayload()
		payload.Name = el.Name
		payload.UniqueId = a.mqtt.GetClientId() + "_" + el.UniqueId
		payload.CommandTopic = a.mqtt.GetActionTopic() + el.UniqueId

		if err := a.mqtt.PublishDiscoveryPayload(payload, "button"); err != nil {
			fmt.Println(err)
		}
		go a.mqtt.Subscribe(payload.CommandTopic, el.Action)
	}
}

// The function updates the states for all the registered sensors
func (a *Agent) updateSensorsState() {
	for _, el := range a.sensorCollectors {
		payload, err := el.CollectorFunc()
		if err != nil {
			fmt.Println(err)
			continue
		}
		a.mqtt.PublishStateUpdate(el.UniqueId, payload)
	}
}

// updates the agent availability
func (a *Agent) updateAvailability(available bool) {
	a.mqtt.PublishAvailabilityPayload(available)
}

