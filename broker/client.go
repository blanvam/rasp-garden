package broker

import "context"

// CallbackHandler interface func definition for a broker callback handler
type CallbackHandler func(c context.Context, topic string, clientID string, config []byte)

// CredentialsProvider should return the current username and password for the MQTT client to use.
type CredentialsProvider func() (username string, password string)

// Client interface definition for a mqtt client conexion
type Client interface {
	IsConnected(c chan bool)
	Connect(c chan error)
	Disconnect(c chan error, quiesce uint)
	Publish(c chan error, topic string, qos uint8, payload interface{})
	Subscribe(c chan error, topic string, qos uint8, callback CallbackHandler)
	Unsubscribe(c chan error, topic string)
}
