package models

import "time"

type DataPoint struct {
	CreatedAt time.Time
	Field     string
	Value     float64
}

type DefaultBrokerPayload struct {
	UnixTimeStamp int64 `json:"unixTimeStamp"`
	Data          []struct {
		Field string  `json:"field"`
		Value float64 `json:"value"`
	} `json:"data"`
}

type Subscriber interface {
	Subscribe(Config, chan<- DataPoint) error
	Finish() error
}

type Flusher interface {
	Connect(Config, <-chan DataPoint) error
	Finish() error
}
