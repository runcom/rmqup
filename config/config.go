package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type OldSoundRabbitMQ struct {
	RabbitMQConfiguration `yaml:"old_sound_rabbit_mq"`
}

type RabbitMQConfiguration struct {
	Connections map[string]Connection `yaml:"connections"`
	Producers   map[string]Producer   `yaml:"producers"`
	Consumers   map[string]Consumer   `yaml:"consumers"`
}

type Connection struct {
	Host     string `yaml:"host,omitempty"`
	Lazy     string `yaml:"lazy,omitempty"`
	Port     string `yaml:"port,omitempty"`
	User     string `yaml:"user,omitempty"`
	Password string `yaml:"password,omitempty"`
	Vhost    string `yaml:"vhost,omitempty"`
}

type Producer struct {
	Connection      string       `yaml:"connection"`
	ExchangeOptions ExchangeOpts `yaml:"exchange_options"`
}

type Consumer struct {
	Connection      string       `yaml:"connection"`
	ExchangeOptions ExchangeOpts `yaml:"exchange_options"`
	QueueOptions    QueueOpts    `yaml:"queue_options"`
	// service to launch inside php command with the json body
	Callback string `yaml:"callback"`
}

type ExchangeOpts struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type QueueOpts struct {
	Name        string              `yaml:"name"`
	RoutingKeys []string            `yaml:"routing_keys"`
	Args        map[string][]string `yaml:"arguments"`
}

func (config *OldSoundRabbitMQ) Parse(path string) error {
	configFile, err := getConfigFile(path)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(configFile, config); err != nil {
		return err
	}

	// config real validation here, just an example..
	// could set a default if empty for example.. etc etc
	//if c.Configuration.Url == "" {
	//return errors.New("rabbitmq_url cannot be empty!")
	//}

	return nil
}

func getConfigFile(path string) ([]byte, error) {
	if _, err := os.Stat(path); err != nil {
		switch {
		case os.IsNotExist(err):
			return nil, fmt.Errorf("Configuration file does not exist at %s", path)
		default:
			return nil, err
		}
	}

	absConfig, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadFile(absConfig)
	if err != nil {
		return nil, err
	}

	return content, nil
}
