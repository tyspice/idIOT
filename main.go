package main

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/tyspice/idIOT/adapters"
	"github.com/tyspice/idIOT/models"
	"gopkg.in/yaml.v2"
)

func main() {
	const (
		etcPath  = "/etc/idIOT/config.yaml"
		homePath = ".config/idIOT/config.yaml"
	)

	var (
		cfg         models.Config
		wg          sync.WaitGroup
		cfgFilePath string
		buffer      = newDataBuffer()
		receiveChan = make(chan models.DataPoint, 10)
		flushChan   = make(chan []models.DataPoint, 5)
	)

	// Init config
	_, err := os.Stat(etcPath)
	etcFileExists := !errors.Is(err, os.ErrNotExist)

	if etcFileExists {
		cfgFilePath = etcPath
	} else {
		homeDir, err := os.UserHomeDir()
		check(err)
		cfgFilePath = filepath.Join(homeDir, homePath)
	}

	content, err := os.ReadFile(cfgFilePath)
	check(err)

	err = yaml.Unmarshal(content, &cfg)
	check(err)

	// Set up broker adapter and start pushing received
	// data points into the buffer
	go bufferFromChannel(buffer, receiveChan, &wg)
	broker := adapters.NewBroker(&cfg)
	broker.Subscribe(cfg, receiveChan)

	// Set up database adapter
	db := adapters.NewDB(&cfg)
	db.Connect(cfg, flushChan)
	go flushAtInterval(buffer, flushChan, time.Duration(cfg.FlushInterval)*time.Second)

	wg.Wait()
}

// UTIL
func check(err error) {
	if err != nil {
		panic(err)
	}
}

// DATA BUFFER
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

// GOROUTINES
func flushAtInterval(buffer *dataBuffer, out chan<- []models.DataPoint, duration time.Duration) {
	for {
		time.Sleep(duration)
		buffer.flushToChannel(out)
	}
}

func bufferFromChannel(buffer *dataBuffer, in <-chan models.DataPoint, wg *sync.WaitGroup) {
	wg.Add(1)
	for {
		dp, ok := <-in
		if !ok {
			break
		}
		buffer.enqueue(dp)
	}
	wg.Done()
}
