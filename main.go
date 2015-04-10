package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/ogier/pflag"
	"github.com/runcom/rmqup/config"
	"github.com/runcom/rmqup/consumer"
)

var (
	flConfig     = pflag.StringP("config", "c", "./app/config/bundles/rabbitmq.yml", "RabbitMQBundle configuration")
	flSfPath     = pflag.StringP("sfpath", "p", "./", "Symfony2 application path")
	flWorker     = pflag.StringP("consumer", "w", "", "Consumer name to start")
	flConnection = pflag.String("connection", "default", "Connection name to use")
	// connection params
)

func main() {
	pflag.Parse()

	var c config.OldSoundRabbitMQ
	if err := c.Parse(*flConfig); err != nil {
		logrus.Fatal(err)
		// exit
		os.Exit(1)
	}

	if *flWorker == "" {
		logrus.Fatal("Please provide a consumer name")
		// exit
		os.Exit(1)
	}

	consumer, err := consumer.NewConsumer(*flWorker, *flConnection, &c)
	if err != nil {
		panic(err)
	}

	if err := consumer.Connect(); err != nil {
		panic(err)
	}
	defer consumer.Disconnect()

	forever := make(chan struct{})
	<-forever
}
