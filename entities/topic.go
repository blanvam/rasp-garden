package entities

import (
	"github.com/eclipse/paho.mqtt.golang"
)

// MQTT is a entity which holds information about the mqtt client
type Topic struct {
	id       string
	qos      int
	client   *mqtt.Client
	options  *mqtt.ClientOptions
	token    string
	mesagges []Message
}
