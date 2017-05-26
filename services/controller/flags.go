package main

import (
	"flag"
)

const (
	flagMongoHost    = "mongohost"
	defaultMongoHost = "localhost"

	flagRabbitHost    = "rabbithost"
	//defaultRabbitHost = "amqp://guest:guest@127.0.0.1:5672/"
		defaultRabbitHost = "amqp://tgdcxdtv:gb0Robr4fcp2@m20.cloudmqtt.com:13595/"
	

)

var (
	confMongoHost  = flag.String(flagMongoHost, defaultMongoHost, "mongo host address, ip:port")
	confRabbitHost = flag.String(flagRabbitHost, defaultRabbitHost, "rabbitmq host address, amqp://user:password@ip:port/")
)
