package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/tyspice/idIOT/adapters"
	"github.com/tyspice/idIOT/models"
	"gopkg.in/yaml.v2"
)

// Data Buffer
type dataBuffer struct {
	mu   sync.Mutex
	data []models.DataPoint
}

func newDataBuffer() *dataBuffer {
	return &dataBuffer{}
}

func (b *dataBuffer) enqueue(dp models.DataPoint) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.data = append(b.data, dp)
}

func (b *dataBuffer) flushToChannel(out chan<- []models.DataPoint) {
	b.mu.Lock()
	defer b.mu.Unlock()
	out <- b.data
	b.data = b.data[:0]
}

func main() {
	var cfg models.Config

	homeDir, err := os.UserHomeDir()
	check(err)

	configPath := filepath.Join(homeDir, ".config/idIOT/config.yaml")

	content, err := os.ReadFile(configPath)
	check(err)

	err = yaml.Unmarshal(content, &cfg)
	check(err)

	buffer := newDataBuffer()
	receiveChan := make(chan models.DataPoint, 10)
	flushChan := make(chan []models.DataPoint, 5)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		for {
			dp, ok := <-receiveChan
			if !ok {
				break
			}
			buffer.enqueue(dp)
		}
		wg.Done()
	}()

	go func() {
		for {
			time.Sleep(time.Duration(cfg.FlushInterval) * time.Second)
			buffer.flushToChannel(flushChan)
		}
	}()

	go func() {
		for {
			flush := <-flushChan
			for _, dp := range flush {
				fmt.Printf("%+v\n", dp)
			}
		}
	}()

	broker := adapters.NewBroker(&cfg)
	broker.Subscribe(cfg, receiveChan)
	wg.Wait()
}

// for development
func check(err error) {
	if err != nil {
		panic(err)
	}
}
