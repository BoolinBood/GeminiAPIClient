package mqtt

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
)

func PublishAudio(topic string, fileName string) error {
	// Read the .wav file into a byte slice
	audioData, err := os.ReadFile(fileName)
	if err != nil {
		return fmt.Errorf("failed to read audio file: %v", err)
	}

	opts := GetClientOptions("audio_client")

	// Create and connect the MQTT client
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %v", token.Error())
	}

	// Publish the .wav file data as the payload to the specified topic
	if token := client.Publish(topic, 0, false, audioData); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish message: %v", token.Error())
	}

	// Disconnect the client
	client.Disconnect(250)

	return nil
}
