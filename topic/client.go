package topic

type CallbackHandler func(topic string, clientID string, config []byte)

// Client interface definition for a mqtt conexion
type Client interface {
	IsConnected(c chan bool)
	Connect(c chan error)
	Disconnect(c chan error, quiesce uint)
	Publish(c chan error, topic string, qos uint8, payload interface{})
	Subscribe(c chan error, topic string, qos uint8, callback CallbackHandler)
	Unsubscribe(c chan error, topic string)
}
