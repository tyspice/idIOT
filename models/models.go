package models

import "time"

type DataPoint struct {
	CreatedAt time.Time
	Field     string
	Value     float64
	Unit      string
}

type DefaultBrokerPayload struct {
	Timestamp int64  `json:"timestamp"`
	SensorId  string `json:"sensorId"`
	Data      []struct {
		Field string  `json:"field"`
		Value float64 `json:"value"`
		Unit  string  `json:"unit"`
	} `json:"data"`
}

type Subscriber interface {
	Subscribe(Config, chan<- DataPoint) error
	Finish() error
}

type Flusher interface {
	Connect(Config, <-chan []DataPoint) error
	Finish() error
}
