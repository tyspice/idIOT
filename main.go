package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/tyspice/idIOT/adapters"
	"github.com/tyspice/idIOT/models"
	"gopkg.in/yaml.v2"
)

// for development
func check(err error) {
	if err != nil {
		panic(err)
	}
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

	receiveChan := make(chan models.DataPoint, 5)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		for {
			dp, ok := <-receiveChan
			if !ok {
				break
			}
			fmt.Printf("%+v\n", dp)
		}
		wg.Done()
	}()

	broker := adapters.NewBroker(&cfg)
	broker.Subscribe(cfg, receiveChan)
	wg.Wait()
}
