package adapters

import (
	"github.com/tyspice/idIOT/adapters/mosquitto"
	"github.com/tyspice/idIOT/models"
)

var brokerMap = map[string]func() models.Subscriber{
	"mosquitto": mosquitto.New,
}

func NewBroker(cfg *models.Config) models.Subscriber {
	return brokerMap[cfg.Broker.Adapter]()
}
