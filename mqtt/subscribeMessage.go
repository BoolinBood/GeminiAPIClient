package mqtt

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// SubscribeMessage connects to an MQTT server, subscribes to a topic, and handles incoming messages.
func SubscribeMessage(broker string, port int, topic string, clientID string) error {
	// Create MQTT client options
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(clientID)
	opts.SetCleanSession(true)

	// Create the MQTT client
	client := mqtt.NewClient(opts)

	// Define the message handler
	messageHandler := func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
	}

	// Set up message handler for received messages
	opts.SetDefaultPublishHandler(messageHandler)

	// Connect to the broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	defer client.Disconnect(250) // Ensure disconnection after done

	// Subscribe to the topic
	token := client.Subscribe(topic, 0, nil) // QoS=0, no message handler (we have one globally)
	token.Wait()                             // Wait until subscription is complete

	// Keep the program running to listen for messages
	for {
		time.Sleep(1 * time.Second)
	}

	return nil
}
