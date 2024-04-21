package mosquitto

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
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
	addr := fmt.Sprintf("tcp://%s:%d", cfg.Broker.IpAddr, cfg.Broker.Port)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(addr)
	opts.SetUsername(cfg.Broker.Username)
	opts.SetPassword(cfg.Broker.Password)
	opts.OnConnect = connectHandler
	opts.SetDefaultPublishHandler(m.handleMessage)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	token := client.Subscribe("topic/test", 1, nil)
	token.Wait()
	return nil
}

func (m *MosquittoClient) Finish() error {
	m.alive = false
	close(m.out)
	return nil
}

func (m *MosquittoClient) handleMessage(_ mqtt.Client, msg mqtt.Message) {
	dp := models.DataPoint{
		CreatedAt: time.Now(),
		Field:     string(msg.Payload()),
		Value:     5,
	}
	m.out <- dp
}

var connectHandler mqtt.OnConnectHandler = func(_ mqtt.Client) {
	fmt.Println("MQTT Connected")
}
