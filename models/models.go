package models

import "time"

type DataPoint struct {
	CreatedAt time.Time
	Field     string
	Value     float64
}

type Subscriber interface {
	Subscribe(Config, <-chan DataPoint) error
	Finish() error
}

type Flusher interface {
	Connect(Config, chan<- DataPoint) error
	Finish() error
}
