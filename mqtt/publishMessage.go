package mqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	TextToBeSpeechTopic = "action/audio/speech"
	GifKeywordTopic     = "ai/gif_keyword"
)

func PublishMessage(topic string, message string) error {
	// Create the MQTT client options

	opts := GetClientOptions("text_client")

	// Create the client
	client := mqtt.NewClient(opts)

	// Connect to the broker
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	defer client.Disconnect(250) // Disconnect when done

	// Publish the message
	token := client.Publish(topic, 2, false, message) // QoS=0, non-persistent
	token.Wait()                                      // Wait for the publish token to complete

	return token.Error()
}
