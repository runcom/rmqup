package consumer

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/runcom/rmqup/config"
	"github.com/streadway/amqp"
)

type Consumer struct {
	ready         bool
	name          string
	connectionUrl string
	connection    *amqp.Connection
}

func (c *Consumer) Connect() error {
	if !c.ready {
		return fmt.Errorf("Can't use non configured consumer %s", c.name)
	}

	conn, err := amqp.DialConfig(c.connectionUrl, amqp.Config{Heartbeat: 60 * time.Second})
	if err != nil {
		return err
	}

	logrus.Infof("Connected to RabbitMQ at %s", c.connectionUrl)
	c.connection = conn

	return nil
}

func NewConsumer(name, connection string, conf *config.OldSoundRabbitMQ) (*Consumer, error) {
	conn, err := conf.GetConnection(connection)
	if err != nil {
		return nil, err
	}
	consumer := &Consumer{
		ready:         true,
		name:          name,
		connectionUrl: buildConnectionUrl(conn),
	}
	//consumerConfig, err := conf.GetConsumer(name)
	//if err != nil {
	//return nil, err
	//}

	return consumer, nil
}

func buildConnectionUrl(conf *config.Connection) string {
	var port string
	if conf.Port != "" {
		port = ":" + conf.Port
	}
	return fmt.Sprintf(
		"amqp://%s:%s@%s%s/%s",
		conf.User,
		conf.Password,
		conf.Host,
		port,
		conf.Vhost,
	)
}
