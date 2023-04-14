package agent

import (
	"github.com/bducha/assistagent/system"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func Shutdown(client MQTT.Client, msg MQTT.Message) {
	system.Shutdown()
}