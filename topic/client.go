package topic

type CallbackHandler func(topic string, clientID string, config []byte)

// Client interface definition for a mqtt conexion
type Client interface {
	IsConnected(c chan bool)
	Disconnect(c chan error, quiesce uint)
	Connect(c chan error)
	Publish(c chan error, topic string, qos uint8, payload interface{})
	Subscribe(c chan error, topic string, qos uint8, callback CallbackHandler)
	Unsubscribe(c chan error, topic string)
}
