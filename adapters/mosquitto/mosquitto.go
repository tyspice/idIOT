package mosquitto

import (
	"encoding/json"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tyspice/idIOT/models"
)

type MosquittoClient struct {
	alive  bool
	out    chan<- models.DataPoint
	client mqtt.Client
}

func New() models.Subscriber {
	return &MosquittoClient{
		alive: false,
	}
}

func (m *MosquittoClient) Subscribe(cfg models.Config, outputChannel chan<- models.DataPoint) error {
	m.alive = true
	m.out = outputChannel
	addr := fmt.Sprintf("tcp://%s:%d", cfg.Broker.IpAddr, cfg.Broker.Port)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(addr)
	opts.SetUsername(cfg.Broker.Username)
	opts.SetPassword(cfg.Broker.Password)
	opts.OnConnect = func(_ mqtt.Client) {
		fmt.Println("MQTT Connected")
	}
	opts.SetDefaultPublishHandler(m.HandleMessage)
	m.client = mqtt.NewClient(opts)
	if token := m.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	for _, topic := range cfg.Broker.Topics {
		token := m.client.Subscribe(topic, 1, nil)
		token.Wait()
	}
	return nil
}

func (m *MosquittoClient) Finish() error {
	m.alive = false
	close(m.out)
	return nil
}

func (m *MosquittoClient) HandleMessage(_ mqtt.Client, msg mqtt.Message) {
	dps, err := defaultPayloadHandler(msg.Payload())
	if err != nil {
		fmt.Println("Error converting payload")
	}
	for _, dp := range dps {
		m.out <- dp
	}
}

func defaultPayloadHandler(payload []byte) ([]models.DataPoint, error) {
	var p models.DefaultBrokerPayload
	var dps []models.DataPoint
	if err := json.Unmarshal(payload, &p); err != nil {
		return nil, err
	}
	createdAt := time.Unix(p.UnixTimeStamp, 0)
	for _, datum := range p.Data {
		dps = append(dps, models.DataPoint{CreatedAt: createdAt, Field: datum.Field, Value: datum.Value})
	}
	return dps, nil
}
