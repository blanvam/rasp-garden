package topic

// Client interface definition for a mqtt conexion
type Client interface {
	IsConnected() bool
	Disconnect(uint)
	Connect(c chan error) error
	Publish(c chan error, topic string, qos uint8, payload interface{})
	Subscribe(c chan []byte, topic string, qos uint8)
	Unsubscribe(c chan error, topic string)
}
