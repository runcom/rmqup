package main

import (
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/ogier/pflag"
	"github.com/runcom/up/config"
)

var (
	flConfig = pflag.StringP("config", "c", "./app/config/bundles/rabbitmq.yml", "RabbitMQBundle configuration")
	flSfPath = pflag.StringP("sfpath", "p", "./", "Symfony2 application path")
	flWorker = pflag.StringP("worker", "w", "", "Consumer name to start")
	// connection params
)

type Worker struct {
}

func newWorker() (*Worker, error) {
	// TODO: implement
	return nil, nil
}

func main() {
	pflag.Parse()

	var c config.OldSoundRabbitMQ
	if err := c.Parse(*flConfig); err != nil {
		logrus.Fatal(err)
		// exit
		os.Exit(1)
	}

	fmt.Printf("%v", c)
}
