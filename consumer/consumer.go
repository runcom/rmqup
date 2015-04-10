package consumer

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/runcom/rmqup/config"
	"github.com/streadway/amqp"
)

type Consumer struct {
	ready           bool
	name            string
	connectionUrl   string
	exchangeOptions config.ExchangeOpts
	queueOptions    config.QueueOpts
	connection      *amqp.Connection
	channel         *amqp.Channel
	queue           *amqp.Queue
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

	c.monitorConnection()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	c.channel = ch

	if err := c.channel.ExchangeDeclare(
		c.exchangeOptions.Name,
		c.exchangeOptions.Type,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	queueArgs := amqp.Table{}
	for i, v := range c.queueOptions.Args {
		queueArgs[i] = v[1]
	}

	if err := queueArgs.Validate(); err != nil {
		return err
	}

	q, err := c.channel.QueueDeclare(
		c.queueOptions.Name,
		false,
		true,
		false,
		false,
		queueArgs,
	)
	c.queue = &q

	for _, k := range c.queueOptions.RoutingKeys {
		if err := c.channel.QueueBind(c.queue.Name, k, c.exchangeOptions.Name, false, nil); err != nil {
			logrus.Fatal(err)
			// don't die here?
		}
	}

	return nil
}

func (c *Consumer) Disconnect() error {
	c.channel.Close()

	return c.connection.Close()
}

func (c *Consumer) monitorConnection() {
	errorsConn := c.connection.NotifyClose(make(chan *amqp.Error))

	go func() {
		for err := range errorsConn {
			logrus.Error(err)
			c.Disconnect()
			select {
			case <-time.After(5 * time.Second):
				// just a timeout to avoid flodding rabbitmq
				// implement connection attempts?
			}
			c.Connect()
		}
	}()
}

func NewConsumer(name, connection string, conf *config.OldSoundRabbitMQ) (*Consumer, error) {
	conn, err := conf.GetConnection(connection)
	if err != nil {
		return nil, err
	}

	consumerConfig, err := conf.GetConsumer(name)
	if err != nil {
		return nil, err
	}

	consumer := &Consumer{
		ready:           true,
		name:            name,
		connectionUrl:   buildConnectionUrl(conn),
		exchangeOptions: consumerConfig.ExchangeOptions,
		queueOptions:    consumerConfig.QueueOptions,
	}

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
