package entities

import "time"

// Interval is a entity which holds information about the resources executuib
type Interval struct {
	Name        string
	Description string
	StartAt     time.Time
	StopAt      time.Time
	duration    time.Duration
	repetitions int
}
