package entities

import "time"

// Message is a entity which holds information about a topic message
type Message struct {
	id          string
	content     string
	DeliveredAt time.Time
	ReceivedAt  time.Time
}
