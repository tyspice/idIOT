package victoriaMetrics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/tyspice/idIOT/models"
)

type VMMetric struct {
	Metric struct {
		Name string `json:"__name__"`
		Unit string `json:"unit"`
	} `json:"metric"`
	Values     []float64 `json:"values"`
	Timestamps []int64   `json:"timestamps"`
}

type VMClient struct {
	alive bool
	in    <-chan []models.DataPoint
	addr  string
}

func New() models.Flusher {
	return &VMClient{
		alive: false,
	}
}

func (vm *VMClient) Connect(cfg models.Config, in <-chan []models.DataPoint) error {
	vm.alive = true
	vm.in = in
	vm.addr = "http://" + cfg.DB.IpAddr + ":" + strconv.Itoa(cfg.DB.Port) + "/api/v1/import"
	go func() {
		for dps := range vm.in {
			if !vm.alive {
				break
			}
			go vm.pushToVM(dps)
		}
	}()
	return nil
}

func (vm *VMClient) Finish() error {
	vm.alive = false
	return nil
}

func (vm *VMClient) pushToVM(dps []models.DataPoint) {
	var body []VMMetric
	for _, dp := range dps {
		vmDP := VMMetric{
			Metric: struct {
				Name string `json:"__name__"`
				Unit string `json:"unit"`
			}{
				Name: dp.Field,
				Unit: dp.Unit,
			},
			Values:     []float64{dp.Value},
			Timestamps: []int64{dp.CreatedAt.Unix()},
		}
		body = append(body, vmDP)
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	resp, err := http.Post(vm.addr, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println("Error code: ", resp.StatusCode)
		fmt.Println("Error response from server:", string(bodyBytes))
	}
}
