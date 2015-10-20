package mattress

import (
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"log"
)

type MqttClient struct {
	client *MQTT.Client
}

var Mqtt *MqttClient

func NewMqttClient() *MqttClient {
	clientoptions := MQTT.NewClientOptions().SetClientID("publisher").AddBroker("tcp://127.0.0.1:1883")
	return &MqttClient{
		client: MQTT.NewClient(clientoptions),
	}
}

func (m *MqttClient) Connection() {
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func (m *MqttClient) Publish(topic string, message string) {
	if token := m.client.Publish(topic, 1, false, message); token.Wait() && token.Error() != nil {
		log.Println(token.Error())
		log.Println("Failed to send message")
	}
}

func SetMqttClient(m *MqttClient) {
	Mqtt = m
}

func GetMqttClient() *MqttClient {
	return Mqtt
}
