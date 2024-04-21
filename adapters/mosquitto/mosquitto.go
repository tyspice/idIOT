package mosquitto

import (
	"time"

	"github.com/tyspice/idIOT/models"
)

type MosquittoClient struct {
	alive bool
	out   chan<- models.DataPoint
}

func New() models.Subscriber {
	return &MosquittoClient{
		alive: false,
	}
}

func (m *MosquittoClient) Subscribe(cfg models.Config, outputChannel chan<- models.DataPoint) error {
	m.alive = true
	m.out = outputChannel
	go func() {
		for m.alive {
			m.out <- models.DataPoint{
				CreatedAt: time.Now(),
				Field:     "my test field",
				Value:     45.5,
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return nil
}

func (m *MosquittoClient) Finish() error {
	m.alive = false
	close(m.out)
	return nil
}
