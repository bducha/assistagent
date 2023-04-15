package mqtt

import (
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	mqttClient MQTT.Client
	clientId   string
}

type DiscoveryPayload struct {
	Availability              *DiscoveryPayloadAvailability `json:"availability,omitempty"`
	Device                    *DiscoveryPayloadDevice       `json:"device,omitempty"`
	Name                      string                        `json:"name,omitempty"`
	StateTopic                string                        `json:"state_topic,omitempty"`
	CommandTopic              string                        `json:"command_topic,omitempty"`
	UniqueId                  string                        `json:"unique_id,omitempty"`
	DeviceClass               string                        `json:"device_class,omitempty"`
	ValueTemplate             string                        `json:"value_template,omitempty"`
	UnitOfMeasurement         string                        `json:"unit_of_measurement,omitempty"`
}

type DiscoveryPayloadDevice struct {
	Identifiers []string `json:"identifiers,omitempty"`
	Name        string   `json:"name,omitempty"`
}

type DiscoveryPayloadAvailability struct {
	Topic               string `json:"topic,omitempty"`
	PayloadAvailable    string `json:"payload_available,omitempty"`
	PayloadNotAvailable string `json:"payload_not_available,omitempty"`
}

// temp dev consts

const (
	BROKER_HOST            = "mqtt://10.0.1.1"
	CLIENT_ID              = "assistagent"
	BASE_TOPIC             = "assistagent/"
	STATES_SUB_TOPIC       = "state/"
	ACTION_SUB_TOPIC       = "action/"
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
		clientId:   CLIENT_ID,
	}
}

func (c *Client)Subscribe(topic string, callback MQTT.MessageHandler) {
    // Subscribe to the topic with the provided callback function
    token := c.mqttClient.Subscribe(topic, 0, callback)
    if token.Wait() && token.Error() != nil {
        panic(token.Error())
    }

    // Wait for an interrupt signal to terminate the subscription
    sig := make(chan os.Signal, 1)
    signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
    <-sig

    // Unsubscribe from the topic and disconnect the client
    token = c.mqttClient.Unsubscribe(topic)
    if token.Wait() && token.Error() != nil {
        panic(token.Error())
    }
}

// Get the MQTT client id
func (c *Client) GetClientId() string {
	return c.clientId
}

// Get the base discovery payload for the device for home assistant
func (c *Client) GetBaseDiscoveryPayload() DiscoveryPayload {
	payload := DiscoveryPayload{}
	payload.Availability = &DiscoveryPayloadAvailability{
		Topic:               BASE_TOPIC + AVAILABILITY_SUB_TOPIC,
		PayloadAvailable:    "on",
		PayloadNotAvailable: "off",
	}
	payload.Device = &DiscoveryPayloadDevice{
		Identifiers: []string{c.GetClientId()},
		Name:        c.GetClientId(),
	}
	return payload
}

// Publishes a message on the provided topic
func (c *Client) Publish(topic string, payload string, retain bool) {

	// fmt.Println("publishing on topic", topic)
	// fmt.Println("payload")
	// fmt.Println(payload)
	token := c.mqttClient.Publish(topic, 0, retain, payload)
	token.Wait()
}

// Sends the availability of the agent
func (c *Client) PublishAvailabilityPayload(available bool) {
	payload := "off"
	if available {
		payload = "on"
	}
	c.Publish(c.GetBaseTopic()+AVAILABILITY_SUB_TOPIC, payload, false)
}

// Publishes a discovery payload to home assistant
func (c *Client) PublishDiscoveryPayload(payload DiscoveryPayload, component string) error {
	topic := "homeassistant/" + component + "/" + payload.UniqueId + "/config"

	data, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	c.Publish(topic, string(data), true)
	return nil
}

// Returns the base topic
func (c *Client) GetBaseTopic() string {
	return BASE_TOPIC
}

// Returns the base state topic
func (c *Client) GetStateTopic() string {
	return c.GetBaseTopic() + STATES_SUB_TOPIC
}

// Returns the base action topic
func (c *Client) GetActionTopic() string {
	return c.GetBaseTopic() + ACTION_SUB_TOPIC
}

// Publishes a state update
func (c *Client) PublishStateUpdate(entityId string, payload string) {
	topic := c.GetStateTopic() + entityId
	c.Publish(topic, payload, true)
}
