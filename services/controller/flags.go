package main

import (
	"flag"
)

const (
	flagMongoHost    = "mongohost"
	defaultMongoHost = "localhost"

	flagRabbitHost = "rabbithost"
	//defaultRabbitHost = "amqp://guest:guest@127.0.0.1:5672/"
	defaultRabbitHost = "amqp://xgurluei:hP9F8ElGZHbAQKQRuqfo5jdT2tpqHuZH@puma.rmq.cloudamqp.com/xgurluei"
)

var (
	confMongoHost  = flag.String(flagMongoHost, defaultMongoHost, "mongo host address, ip:port")
	confRabbitHost = flag.String(flagRabbitHost, defaultRabbitHost, "rabbitmq host address, amqp://user:password@ip:port/")
)
