package models

type Config struct {
	FlushInterval int `yaml:"flushInterval"`
	Broker        struct {
		Adapter string   `yaml:"adapter"`
		Uri     string   `yaml:"uri"`
		IpAddr  string   `yaml:"ipAddr"`
		Port    int      `yaml:"port"`
		Topics  []string `yaml:"topics"`
	} `yaml:"broker"`
	DB struct {
		Adapter string `yaml:"adapter"`
		Uri     string `yaml:"uri"`
		IpAddr  string `yaml:"ipAddr"`
		Port    int    `yaml:"port"`
	} `yaml:"db"`
}
