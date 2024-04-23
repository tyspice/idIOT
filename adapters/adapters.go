package adapters

import (
	"github.com/tyspice/idIOT/adapters/mosquitto"
	"github.com/tyspice/idIOT/adapters/victoriaMetrics"
	"github.com/tyspice/idIOT/models"
)

var brokerMap = map[string]func() models.Subscriber{
	"mosquitto": mosquitto.New,
}

var dbMap = map[string]func() models.Flusher{
	"victoriaMetrics": victoriaMetrics.New,
}

func NewBroker(cfg *models.Config) models.Subscriber {
	return brokerMap[cfg.Broker.Adapter]()
}

func NewDB(cfg *models.Config) models.Flusher {
	return dbMap[cfg.DB.Adapter]()
}
