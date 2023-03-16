package mqtt

import (
	"encoding/json"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	mqttClient MQTT.Client
	clientId   string
}

type DiscoveryPayload struct {
	Availability struct {
		Topic               string `json:"topic,omitempty"`
		PayloadAvailable    string `json:"payload_available,omitempty"`
		PayloadNotAvailable string `json:"payload_not_available,omitempty"`
	} `json:"availability,omitempty"`
	Device struct {
		Identifiers []string `json:"identifiers,omitempty"`
		Name        string   `json:"name,omitempty"`
	} `json:"device,omitempty"`
	Name       string `json:"name,omitempty"`
	StateTopic string `json:"state_topic,omitempty"`
	UniqueId   string `json:"unique_id,omitempty"`
}

// temp dev consts

const (
	BROKER_HOST = "mqtt://10.0.1.1"
	CLIENT_ID   = "assistagent"
	BASE_TOPIC = "/assistagent/"
	STATES_SUB_TOPIC = "state/"
	AVAILABILITY_SUB_TOPIC = "health"
)

// Initialize a new MQTT client
func NewClient() Client {
	opts := MQTT.NewClientOptions()
	opts.AddBroker("mqtt://10.0.1.1:1883")
	opts.SetClientID(CLIENT_ID)

	// init and connect the client
	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return Client{
		mqttClient: client,
	}
}

// Get the MQTT client id
func (c *Client) GetClientId() string {
	return c.clientId
}

// Get the base discovery payload for the device for home assistant
func (c *Client) GetBaseDiscoveryPayload() DiscoveryPayload {
	payload := DiscoveryPayload{}
	payload.Availability.Topic = BASE_TOPIC + AVAILABILITY_SUB_TOPIC
	payload.Availability.PayloadAvailable = "on"
	payload.Availability.PayloadNotAvailable = "off"
	payload.Device.Identifiers = []string{c.GetClientId()}
	payload.Device.Name = c.GetClientId()
	return payload
}

// Publishes a message on the provided topic
func (c *Client) Publish(topic string, payload string, retain bool) {

	token := c.mqttClient.Publish(topic, 0, retain, payload)
	token.Wait()
}

// Publishes a discovery payload to home assistant
func (c *Client) PublishDiscoveryPayload(payload DiscoveryPayload, component string) error {
	topic := "/homeassistant/" + component + "/" + payload.UniqueId + "/config"

	data, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	c.Publish(topic, string(data), true)
	return nil
}

func (c *Client) GetBaseTopic() string {
	return BASE_TOPIC
}

func (c *Client) GetStateTopic() string {
	return c.GetBaseTopic() + STATES_SUB_TOPIC
}

// Publishes a state update
func (c *Client) PublishStateUpdate(entityId string, payload string) {
	topic := c.GetStateTopic() + entityId
	c.Publish(topic, payload, true)
}
