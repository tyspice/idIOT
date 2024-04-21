package main

import (
	"fmt"
	"os"
	"path/filepath"

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

	fmt.Printf("Config: %+v\n", cfg)

}
