package victoriaMetrics

import (
	"fmt"

	"github.com/tyspice/idIOT/models"
)

type VMClient struct {
	alive bool
	in    <-chan []models.DataPoint
}

func New() models.Flusher {
	return &VMClient{
		alive: false,
	}
}

func (vm *VMClient) Connect(cfg models.Config, in <-chan []models.DataPoint) error {
	vm.alive = true
	vm.in = in
	go func() {
		for dps := range vm.in {
			if !vm.alive {
				break
			}
			for _, dp := range dps {
				fmt.Printf("%+v\n", dp)
			}
		}
	}()
	return nil
}

func (vm *VMClient) Finish() error {
	vm.alive = false
	return nil
}
