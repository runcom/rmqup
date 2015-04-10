package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/ogier/pflag"
	"github.com/runcom/rmqup/config"
	"github.com/runcom/rmqup/consumer"
)

var (
	flConfig = pflag.StringP("config", "c", "./app/config/bundles/rabbitmq.yml", "RabbitMQBundle configuration")
	flSfPath = pflag.StringP("sfpath", "p", "./", "Symfony2 application path")
	flWorker = pflag.StringP("worker", "w", "", "Consumer name to start")
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

	//fmt.Printf("%v", c)

	consumer, err := consumer.NewConsumer("mailer", "default", &c)
	if err != nil {
		panic(err)
	}
	if err := consumer.Connect(); err != nil {
		panic(err)
	}

	forever := make(chan struct{})
	<-forever
	fmt.Printf("%v", consumer)
}
