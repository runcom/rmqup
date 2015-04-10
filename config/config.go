package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// Root yaml configuration. An anonymous field is used here to get rid of the
// ridondance of having two config struct nested.
type OldSoundRabbitMQ struct {
	rabbitMQConfiguration `yaml:"old_sound_rabbit_mq"`
}

func (c *OldSoundRabbitMQ) GetConnection(name string) (*Connection, error) {
	conn, ok := c.Connections[name]
	if !ok {
		return nil, fmt.Errorf("No connection named %s", name)
	}
	return &conn, nil
}

func (c *OldSoundRabbitMQ) GetConsumer(name string) (*Consumer, error) {
	consumer, ok := c.Consumers[name]
	if !ok {
		return nil, fmt.Errorf("No consumer named %s", name)
	}
	return &consumer, nil
}

// func (c *OldSoundRabbitMQ) GetProducer(name string) *Producer {}

// Main RabbitMQBundle configuration tree
type rabbitMQConfiguration struct {
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
	ConnectionName  string       `yaml:"connection"`
	ExchangeOptions ExchangeOpts `yaml:"exchange_options"`
}

type Consumer struct {
	ConnectionName  string       `yaml:"connection"`
	ExchangeOptions ExchangeOpts `yaml:"exchange_options"`
	QueueOptions    QueueOpts    `yaml:"queue_options"`
	Callback        string       `yaml:"callback"`
	// callback above is the service to launch inside php command with the json body
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

// Parse is used to populate a config.OldSoundRabbitMQ struct from a fs path
// pointing to a rabbitmq.yml file
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

// helper function that validates path and read the file
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
