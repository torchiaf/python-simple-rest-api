package models

type RabbitMQ struct {
	Host     string
	Port     string
	Username string
	Password string
}

type Config struct {
	IsDev    bool
	RabbitMQ RabbitMQ
	Modules  map[string]Module
}

type Module struct {
	Name       string `yaml:"name"`
	NodeName   string `yaml:"nodeName"`
	Type       string `yaml:"type"`
	RoutingKey string `yaml:"routingKey"`
	Devices    []Device
}

type Device struct {
	Name   string `yaml:"name"`
	Type   string `yaml:"type"`
	Config []DeviceConfig
}

type DeviceConfig struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Message struct {
	// TODO refactoring
	Device string `json:"device"`
	Action string `json:"action"`
	// Args   map[string]interface{} `json:"args"`
}
