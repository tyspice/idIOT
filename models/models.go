package models

import "time"

type DataPoint struct {
	CreatedAt time.Time
	Field     string
	Value     float64
}

type BrokerClient struct {
	Broker    string
	Subscribe func() <-chan DataPoint
	Finish    func()
}

type DBClient struct {
	DB      string
	Connect func() chan<- DataPoint
	Finish  func()
}
