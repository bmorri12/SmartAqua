package main

import (
	"flag"
)

const (
	flagRabbitHost    = "rabbithost"
	defaultRabbitHost = "amqp://guest:guest@127.0.0.1:5672/"
)

var (
	confRabbitHost = flag.String(flagRabbitHost, defaultRabbitHost, "rabbitmq host address, amqp://user:password@ip:port/")
)
